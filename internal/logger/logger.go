package logger

import (
	"log"
	"os"
)

type Logger struct {
	LogLevel int
	log      *log.Logger
}

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func NewLogger(LogLevel string) *Logger {
	l := &Logger{
		log: log.New(os.Stdout, "", log.LstdFlags),
	}
	switch LogLevel {
	case "debug", "DEBUG", "Debug":
		l.LogLevel = DEBUG
	case "info", "INFO", "Info":
		l.LogLevel = INFO
	case "warning", "WARNING", "Warning":
		l.LogLevel = WARNING
	case "error", "ERROR", "Error":
		l.LogLevel = ERROR
	case "fatal", "FATAL", "Fatal":
		l.LogLevel = FATAL
	default:
		l.LogLevel = INFO
	}
	return l
}

func (l *Logger) Debug(args ...interface{}) {
	if l.LogLevel == DEBUG {
		l.log.Println(args...)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.LogLevel <= INFO {
		l.log.Println(args...)
	}
}

func (l *Logger) Warning(args ...interface{}) {
	if l.LogLevel <= WARNING {
		l.log.Println(args...)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.LogLevel <= ERROR {
		l.log.Println(args...)
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.LogLevel <= FATAL {
		l.log.Println(args...)
	}
}
