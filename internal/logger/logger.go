package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	LogLevel int
	log      *log.Logger
}

const (
	TEST  = 10
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func New(LogLevel string) *Logger {
	l := &Logger{
		log: log.New(os.Stdout, "", log.LstdFlags),
	}
	switch LogLevel {
	case "test", "TEST", "Test":
		l.LogLevel = TEST
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
		mes := []string{"DEBUG: "}
		for _, v := range args {
			mes = append(mes, fmt.Sprint(v, " "))
		}
		l.log.Print(mes)
	}
}

func (l *Logger) Info(args ...interface{}) {
	if l.LogLevel <= INFO {
		mes := []string{"INFO: "}
		for _, v := range args {
			mes = append(mes, fmt.Sprint(v, " "))
		}
		l.log.Print(mes)
	}
}

func (l *Logger) Warning(args ...interface{}) {
	if l.LogLevel <= WARNING {
		mes := []string{"WARNING: "}
		for _, v := range args {
			mes = append(mes, fmt.Sprint(v, " "))
		}
		l.log.Print(mes)
	}
}

func (l *Logger) Error(args ...interface{}) {
	if l.LogLevel <= ERROR {
		mes := []string{"ERROR: "}
		for _, v := range args {
			mes = append(mes, fmt.Sprint(v, " "))
		}
		l.log.Print(mes)
	}
}
func (l *Logger) Fatal(args ...interface{}) {
	if l.LogLevel <= FATAL {
		mes := []string{"FATAL: "}
		for _, v := range args {
			mes = append(mes, fmt.Sprint(v, " "))
		}
		l.log.Print(mes)
	}
}
