package middleware

import (
	infra "nh-be/infra/observability"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		infra.HttpRequestsInFlight.Inc()
		start := time.Now()

		c.Next()

		infra.HttpRequestsInFlight.Dec()
		duration := time.Since(start).Seconds()

		infra.HttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		infra.HttpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}
