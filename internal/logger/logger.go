package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	log *zap.Logger
}

func New(logLevel string) *Logger {
	var level zap.AtomicLevel
	switch logLevel {
	case "debug", "DEBUG":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info", "INFO":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn", "WARNING":
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error", "ERROR":
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	cfg := zap.Config{
		Level:            level,
		Encoding:         "json",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, _ := cfg.Build()
	return &Logger{log: logger}
}

func (l *Logger) LogRequest(req *http.Request, statusCode int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("uri", req.URL.Path),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
		zap.String("timestamp", time.Now().Format(time.RFC3339)),
	}

	if err != nil {
		fields = append(fields, zap.String("error", err.Error()))
		l.log.Error("Request failed", fields...)
		return
	}

	if l.log.Core().Enabled(zap.DebugLevel) {
		for key, values := range req.Header {
			fields = append(fields, zap.String("header_"+key, values[0]))
		}
		fields = append(fields, zap.String("body", "[LOG_BODY_IMPLEMENTATION]"))
	}

	l.log.Info("Request processed", fields...)
}
