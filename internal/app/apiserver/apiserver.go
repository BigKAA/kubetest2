package apiserver

import (
	"io"
	"net/http"

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

	s.logger.Info("Starting API server")
	s.ConfigRouter()
	s.logger.Info("Listen on: http://" + s.config.BindAddr)
	return http.ListenAndServe(s.config.BindAddr, s.router)
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
	s.router.HandleFunc("/hello", s.HandlerHello())
}

// HandlerHello ...
func (s *APIServer) HandlerHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
