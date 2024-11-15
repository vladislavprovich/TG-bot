package logger

import (
	"github.com/sirupsen/logrus"
)

func NewLogger(LVLlog string) *logrus.Logger {
	logger := logrus.New()

	if LVLlog == "prod" {
		logger.SetLevel(logrus.InfoLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
	return logger
}
