package middleware

import (
	"context"
	"net/http"
	"nh-be/internal/constant"
	"nh-be/internal/infra"

	"github.com/gin-gonic/gin"
)

func RequireAuth(sessionStore infra.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("auth_session")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "unauthorized",
				"message": "Unauthorized",
			})
			return
		}

		userID, err := sessionStore.GetUserSession(c.Request.Context(), cookie.Value)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "invalid_user_id",
				"message": "Invalid user ID in session",
			})
			return
		}

		ctx := context.WithValue(c.Request.Context(), constant.CtxUserId, userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
