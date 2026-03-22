package middlwares

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		log.InfoContext(c.Request.Context(), "request started",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("client_ip", c.ClientIP()),
		)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		log.InfoContext(c.Request.Context(), "request completed",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Duration("latency", latency),
		)
	}
}
