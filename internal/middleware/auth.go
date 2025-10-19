package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
  return func(c *gin.Context) {
    sess := sessions.Default(c)
    userID := sess.Get("user_id")
    if userID == nil {
      c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
      return
    }
    c.Next()
  }
}