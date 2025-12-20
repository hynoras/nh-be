package middleware

import (
	"context"
	"net/http"
	"nh-be/constant"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		userIDStr := sess.Get("user_id")
		if userIDStr == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "Unauthorized",
			})
			return
		}

		// Parse user ID string to UUID and set in context
		userID, err := uuid.Parse(userIDStr.(string))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "invalid_user_id",
				"message": "Invalid user ID in session",
			})
			return
		}

		// Set user ID in context for downstream handlers/services
		ctx := context.WithValue(c.Request.Context(), constant.CtxUserId, userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
