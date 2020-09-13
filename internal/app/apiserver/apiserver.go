package apiserver

import (
	"flag"
	"net/http"
	"os"
	"path/filepath"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// APIServer собственно сервер
type APIServer struct {
	config   *Config
	logger   *logrus.Logger
	router   *mux.Router
	restconf *rest.Config // для доступа к API кластера кубернетес
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

// ConfigRest Конфигурация подключения к rest API кластера
func (s *APIServer) ConfigRest() error {
	// сначал надо понять находимся ли мы внутри кластера
	// для этого проконтролируем наличеие файа
	// /var/run/secrets/kubernetes.io/serviceaccount/token
	if fileExists("/var/run/secrets/kubernetes.io/serviceaccount/token") {
		// внутри кластера
		conf, err := rest.InClusterConfig()
		if err != nil {
			s.logger.Error("Can't config connection to API. ", err)
			return err
		}
		s.restconf = conf
	} else {
		// не в кластере
		var kubeconfig *string
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()

		// use the current context in kubeconfig
		conf, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			s.logger.Error(err)
		}
		s.restconf = conf
	}
	return nil
}

// homeDir ...
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// fileExists Проверяет есть ли такой файл
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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
