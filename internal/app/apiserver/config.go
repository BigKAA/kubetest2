package apiserver

import "os"

// Config содержит конфигурационные параметры
type Config struct {
	BindAddr string
	LogLevel string
}

// NewConfig получает конфигурационные параметры из среды окружения.
// если переменной среды окружения нет, подставляет значения по умолчанию.
func NewConfig() *Config {
	bindAddr, exists := os.LookupEnv("BIND_ADDR")
	if !exists {
		bindAddr = "127.0.0.1:8080"
	}

	logLevel, exists := os.LookupEnv("DEFAULT_LOG_LEVEL")
	if !exists {
		logLevel = "debug"
	}

	return &Config{
		BindAddr: bindAddr,
		LogLevel: logLevel,
	}
}
