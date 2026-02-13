package middleware

import (
	"file-parser/internal/config"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitLogger(config *config.LoggerConfig) gin.HandlerFunc {
	logger := logrus.New()

	setLoggerLevel(logger, config)
	setFormatter(logger, config)

	logger.SetOutput(os.Stdout)

	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		entry := logger.WithFields(logrus.Fields{
			"status":  status,
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"ip":      c.ClientIP(),
			"latency": latency.String(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			switch {
			case status >= 500:
				entry.Error("server error")
			case status >= 400:
				entry.Warn("client error")
			default:
				entry.Info("request completed")
			}
		}
	}
}

// Set logger level
func setLoggerLevel(logger *logrus.Logger, conf *config.LoggerConfig) {
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

func setFormatter(logger *logrus.Logger, conf *config.LoggerConfig) {
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
