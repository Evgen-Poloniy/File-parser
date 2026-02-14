package logutil

import (
	"file-parser/internal/config"

	"github.com/sirupsen/logrus"
)

// Set logger level
func SetLevel(logger *logrus.Logger, conf *config.LoggerConfig) {
	switch conf.Level {
	case config.DebugLevel:
		logger.SetLevel(logrus.DebugLevel)
	case config.InfoLevel:
		logger.SetLevel(logrus.InfoLevel)
	case config.WarnLevel:
		logger.SetLevel(logrus.WarnLevel)
	case config.ErrorLevel:
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

// Set logger format
func SetFormatter(logger *logrus.Logger, conf *config.LoggerConfig) {
	switch conf.Format {
	case config.JsonFormat:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	case config.TextFormat:
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}
}
