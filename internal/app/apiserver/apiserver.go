package apiserver

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// APIServer собственно сервер
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

// New создаёт экземпляр сервера
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Start запускает сервер
func (s *APIServer) Start() error {
	if err := s.ConfigLogger(); err != nil {
		return err
	}
	s.ConfigRouter()
	s.logger.Info("Starting API server Listen on: http://", s.config.BindAddr)

	server := &http.Server{
		Addr:    s.config.BindAddr,
		Handler: s.logRequestHandler(s.router),
		// ErrorLog: s.logger,
	}

	return server.ListenAndServe()
}

// ConfigLogger конфигурация логгера
func (s *APIServer) ConfigLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)
	s.logger.SetFormatter(&logrus.JSONFormatter{})

	return nil
}

// ConfigRouter конфигурирует роутер
func (s *APIServer) ConfigRouter() {
	pods := s.router.PathPrefix("/pods").Subrouter()
	pods.HandleFunc("/", s.HandlerPods())
	pods.HandleFunc("/{ns}", s.HandlerPods())
}

// logRequestHandler логирует запросы
func (s *APIServer) logRequestHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(h, w, r)
		s.logger.WithFields(logrus.Fields{
			"method":         r.Method,
			"remote_address": r.RemoteAddr,
			"request_uri":    r.RequestURI,
			"user_agent":     r.UserAgent(),
			"status":         m.Code,
			"bytes":          m.Written,
			"duration":       m.Duration,
		}).Info("request")
	}
}
