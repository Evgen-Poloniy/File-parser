package middleware

import (
	"file-parser/internal/config"
	"file-parser/internal/logutil"
	"strings"

	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func InitLogger(config *config.LoggerConfig) gin.HandlerFunc {
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logutil.SetLevel(logger, config)
	logutil.SetFormatter(logger, config)

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
			errs := strings.Join(c.Errors.Errors(), "; ")
			entry.Error(errs)
			// entry.Error(c.Errors.String())
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
