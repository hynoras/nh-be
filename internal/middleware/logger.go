package middleware

import (
	"log/slog"
	"nh-be/internal/constant"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		hostName, _ := os.Hostname()

		requestID, _ := c.Get(string(constant.CtxRequestId))
		service, _ := c.Get(string(constant.CtxService))

		slog.Info("request",
			"request_id", requestID,
			"host_name", hostName,
			"service", service,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latency.Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func WithService(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(string(constant.CtxService), name)
		c.Next()
	}
}
