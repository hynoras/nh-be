package auth

import (
	"net/http"
	"nh-be/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LoginHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginDto
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.MakeErrorResponse(
				c, 
				http.StatusBadRequest, 
				"Invalid request format", 
				err.Error(),
			)
			return
		}

		user, err := s.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			utils.MakeErrorResponse(
				c, 
				http.StatusUnauthorized, 
				"Invalid email or password", 
				err.Error(),
			)
			return
		}

		sess := sessions.Default(c)
		sess.Set("user_id", user.ID.String())
		if err := sess.Save(); err != nil {
			utils.MakeErrorResponse(
				c, 
				http.StatusInternalServerError, 
				"Session save failed", 
				err.Error(),
			)
			return
		}

		// Get the session ID from the session store

		resp := LoginResponseDto{
			User: UserResponseDto{
				ID: user.ID,
				Username: user.Username,
				Email: user.Email,
				Role: user.Role,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		}
		utils.MakeSuccessResponse(c, "User logged in successfully", resp)	
	}
}

func LogoutHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := s.Logout(c); err != nil {
			utils.MakeErrorResponse(
				c, 
				http.StatusInternalServerError, 
				"Failed to logout", 
				err.Error(),
			)
			return
		}
		utils.MakeSuccessResponse(c, "User logged out successfully", nil)
	}
}
