package middlewares

import (
	"strconv"
	"time"

	"github.com/OnYyon/oregon-api-gateway/pkg/metrics"
	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method
		metrics.InFlightRequests.WithLabelValues(method, path).Inc()
		defer metrics.InFlightRequests.WithLabelValues(method, path).Dec()

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		duration := time.Since(start).Seconds()
		metrics.RequestsTotal.WithLabelValues(method, path, status).Inc()
		metrics.RequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
