package middleware

import (
	"context"
	"net/http"
	"nh-be/internal/constant"
	"nh-be/internal/infra"
	"nh-be/internal/utils/httputil"

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
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to get session",
				"message": err,
			})
			return
		}

		parsedUserId, err := httputil.ParseStringToUUID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Failed to parse user ID",
				"message": err,
			})
			return
		}
		ctx := context.WithValue(c.Request.Context(), constant.CtxUserId, parsedUserId)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
