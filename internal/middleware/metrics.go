package middleware

import (
	obs "nh-be/internal/platform/observability"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		obs.HttpRequestsInFlight.Inc()
		start := time.Now()

		c.Next()

		obs.HttpRequestsInFlight.Dec()
		duration := time.Since(start).Seconds()

		obs.HttpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			strconv.Itoa(c.Writer.Status()),
		).Inc()

		obs.HttpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Observe(duration)
	}
}
