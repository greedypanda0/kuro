package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTP HTTPConfig
	Log  LogConfig
}

type HTTPConfig struct {
	Addr            string
	ShutdownTimeout time.Duration
}

type LogConfig struct {
	Level       string
	Development bool
}

func Load() Config {
	cfg := Config{
		HTTP: HTTPConfig{
			Addr:            ":8080",
			ShutdownTimeout: 10 * time.Second,
		},
		Log: LogConfig{
			Level:       "info",
			Development: false,
		},
	}

	if v := os.Getenv("REMOTE_HTTP_ADDR"); v != "" {
		cfg.HTTP.Addr = v
	}

	if v := os.Getenv("REMOTE_HTTP_SHUTDOWN_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.HTTP.ShutdownTimeout = d
		}
	}

	if v := os.Getenv("REMOTE_LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}

	if v := os.Getenv("REMOTE_LOG_DEV"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			cfg.Log.Development = b
		}
	}

	return cfg
}
