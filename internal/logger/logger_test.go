package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		expected int
	}{
		{"Default to INFO", "", INFO},
		{"Debug Level", "debug", DEBUG},
		{"Info Level", "INFO", INFO},
		{"Warning Level", "warning", WARNING},
		{"Error Level", "ERROR", ERROR},
		{"Fatal Level", "FATAL", FATAL},
		{"Invalid Level defaults to INFO", "invalid", INFO},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.logLevel)
			if logger.LogLevel != tt.expected {
				t.Errorf("expected log level %d, got %d", tt.expected, logger.LogLevel)
			}
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := &Logger{
		LogLevel: DEBUG,
		log:      log.New(&buf, "", 0),
	}

	tests := []struct {
		name     string
		method   func(...interface{})
		level    int
		message  string
		expected bool
	}{
		{"Debug at DEBUG", logger.Debug, DEBUG, "debug message", true},
		{"Info at DEBUG", logger.Info, DEBUG, "info message", true},
		{"Warning at DEBUG", logger.Warning, DEBUG, "warning message", true},
		{"Error at DEBUG", logger.Error, DEBUG, "error message", true},
		{"Fatal at DEBUG", logger.Fatal, DEBUG, "fatal message", true},

		{"Debug at INFO", logger.Debug, INFO, "debug message", false},
		{"Info at INFO", logger.Info, INFO, "info message", true},
		{"Warning at INFO", logger.Warning, INFO, "warning message", true},
		{"Error at INFO", logger.Error, INFO, "error message", true},
		{"Fatal at INFO", logger.Fatal, INFO, "fatal message", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			logger.LogLevel = tt.level
			tt.method(tt.message)
			output := buf.String()
			if tt.expected && !strings.Contains(output, tt.message) {
				t.Errorf("expected message %q in output, got %q", tt.message, output)
			}
			if !tt.expected && strings.Contains(output, tt.message) {
				t.Errorf("did not expect message %q in output, but got %q", tt.message, output)
			}
		})
	}
}
