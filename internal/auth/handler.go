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

		userRes, err := s.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			utils.MakeErrorResponse(
				c,
				http.StatusUnauthorized,
				"Invalid email or password",
				err.Error(),
			)
			return
		}

		// Save user id to session
		sess := sessions.Default(c)
		sess.Set("user_id", userRes.ID.String())
		if err := sess.Save(); err != nil {
			utils.MakeErrorResponse(
				c,
				http.StatusInternalServerError,
				"Session save failed",
				err.Error(),
			)
			return
		}

		resp := LoginResponseDto{
			User: UserResponseDto{
				ID:        userRes.ID,
				Username:  userRes.Username,
				Email:     userRes.Email,
				CreatedAt: userRes.CreatedAt,
				UpdatedAt: userRes.UpdatedAt,
			},
		}
		utils.MakeSuccessResponse(c, http.StatusOK, "User logged in successfully", resp)
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
		utils.MakeSuccessResponse(c, http.StatusOK, "User logged out successfully", nil)
	}
}

func ChangePasswordHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.ParseStringToUUID(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}

		var req ChangePasswordDto
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		err = s.ChangePassword(c.Request.Context(), userID, req)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Failed to update user password", err.Error())
			return
		}
		utils.MakeSuccessResponse(c, http.StatusOK, "User password changed successfully", nil)
	}
}
