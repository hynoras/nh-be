package auth

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginDto
		if err := c.ShouldBindJSON(&req); err != nil {
			c.Set("message", "Invalid request format")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := s.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			c.Set("message", "Invalid email or password")
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		sess := sessions.Default(c)
		sess.Set("user_id", user.ID.String())
		if err := sess.Save(); err != nil {
			c.Set("message", "Session save failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := LoginResponseDto{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		}
		c.Set("message", "User logged in successfully")
		c.JSON(http.StatusOK, resp)
	}
}