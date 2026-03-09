package middleware

import (
	"context"
	"nh-be/internal/constant"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.GetHeader("X-Request-ID")

		if id == "" || len(id) > 128 {
			id = uuid.New().String()
		}

		c.Set(string(constant.CtxRequestId), id)
		c.Writer.Header().Set("X-Request-ID", id)

		ctx := context.WithValue(c.Request.Context(), constant.CtxRequestId, id)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
