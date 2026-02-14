package auth

import (
	"net/http"
	"nh-be/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoginHandler godoc
// @Summary User login
// @Description Authenticate user with email and password, returns user information with permissions and creates a session
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body LoginDto true "Login credentials (email and password)"
// @Success 200 {object} utils.SuccessResponse{data=LoginResponseDto} "Successfully logged in"
// @Failure 400 {object} utils.ErrorResponse "Invalid request format"
// @Failure 401 {object} utils.ErrorResponse "Invalid email or password"
// @Failure 500 {object} utils.ErrorResponse "Session save failed"
// @Router /auth/login [post]
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

		userRes, permRes, err := s.Login(c.Request.Context(), req.Email, req.Password)
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
				ID:          userRes.ID,
				Username:    userRes.Username,
				Email:       userRes.Email,
				Permissions: permRes,
				CreatedAt:   userRes.CreatedAt,
				UpdatedAt:   userRes.UpdatedAt,
			},
		}
		utils.MakeSuccessResponse(c, http.StatusOK, "User logged in successfully", resp)
	}
}

// LogoutHandler godoc
// @Summary User logout
// @Description Clear user session and log out from the system
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse "Successfully logged out"
// @Failure 500 {object} utils.ErrorResponse "Failed to logout"
// @Security SessionAuth
// @Router /auth/logout [post]
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

// ChangePasswordHandler godoc
// @Summary Change user password
// @Description Update user password with new password and confirmation
// @Tags Authentication
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID format)"
// @Param request body ChangePasswordDto true "New password and confirmation"
// @Success 200 {object} utils.SuccessResponse "Password changed successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid user ID or request body"
// @Security SessionAuth
// @Router /auth/users/{id}/change-password [put]
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

func SignUpHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SignUpDto
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		err := s.SignUp(c.Request.Context(), req)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Failed to create user", err.Error())
			return
		}
		utils.MakeSuccessResponse(c, http.StatusOK, "User created successfully", nil)
	}
}
