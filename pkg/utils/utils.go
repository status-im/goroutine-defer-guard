package utils

import (
	"log"
	"os"
	"strconv"
)

func URI(path string, line int) string {
	return path + ":" + strconv.Itoa(line)
}

type Logger struct {
	logger *log.Logger
}

func BuildLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l Logger) Errorf(format string, args ...any) {
	l.logger.Printf("ERROR: "+format, args...)
}

func (l Logger) Warnf(format string, args ...any) {
	l.logger.Printf("WARN: "+format, args...)
}

func (l Logger) Infof(format string, args ...any) {
	l.logger.Printf("INFO: "+format, args...)
}
