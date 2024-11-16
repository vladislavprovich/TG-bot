package logger

import (
	"github.com/sirupsen/logrus"
)

const (
	debug = "debug"
	info  = "info"
)

func NewLogger(cfg ConfigLogger) *logrus.Logger {
	logger := logrus.New()

	switch cfg.Level {
	case debug:
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
	case info:
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.TextFormatter{})
	default:
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger
}
