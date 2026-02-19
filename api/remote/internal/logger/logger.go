package logger

import (
	"strings"
	"time"

	"github.com/greedypanda0/kuro/api/remote/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.Logger

func Init(cfg config.LogConfig) *Logger {
	var zapCfg zap.Config
	if cfg.Development {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}
	zapCfg.Encoding = "console"
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if level := parseLevel(cfg.Level); level != nil {
		zapCfg.Level = zap.NewAtomicLevelAt(*level)
	}

	log, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}

	return log
}

func parseLevel(level string) *zapcore.Level {
	if level == "" {
		return nil
	}

	var lvl zapcore.Level
	if err := lvl.Set(strings.ToLower(level)); err != nil {
		return nil
	}

	return &lvl
}

func String(key, value string) zap.Field {
	return zap.String(key, value)
}

func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

func Duration(key string, value time.Duration) zap.Field {
	return zap.Duration(key, value)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}
