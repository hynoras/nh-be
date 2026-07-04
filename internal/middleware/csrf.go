// middleware/csrf.go
package middleware

import (
	"fmt"
	"net/http"

	"nh-be/internal/utils/crypto"
	"nh-be/internal/utils/httputil"

	"github.com/gin-gonic/gin"
)

func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip safe methods
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Validate: X-CSRF-Token header must match csrf_token cookie
		cookieToken, err := c.Request.Cookie("csrf_token")
		headerToken := c.GetHeader("X-CSRF-Token")

		fmt.Println("cookieToken", cookieToken.Value)
		fmt.Println("headerToken", headerToken)

		if err != nil || cookieToken.Value == "" || cookieToken.Value != headerToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "csrf_validation_failed",
				"message": "CSRF token missing or invalid",
			})
			return
		}
		c.Next()
	}
}

func SetCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Request.Cookie("csrf_token"); err != nil {
			token, err := crypto.GenerateToken()
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error":   "internal_server_error",
					"message": "Failed to generate CSRF token",
				})
				return
			}
			http.SetCookie(c.Writer, httputil.GetCSRFTokenCookie(token))
		}
		c.Next()
	}
}
