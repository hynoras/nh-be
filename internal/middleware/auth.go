package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"nh-be/internal/constant"
	"nh-be/internal/platform/session"
	"nh-be/internal/utils/httputil"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RequireAuth(sessionStore session.SessionStore) gin.HandlerFunc {
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
			if errors.Is(err, redis.Nil) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "session_expired",
					"message": "Session not found or expired",
				})
				return
			}
			slog.Error("redis session lookup failed", "error", err)
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"error":   "service_unavailable",
				"message": "Session service temporarily unavailable",
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
