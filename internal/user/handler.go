package user

import (
	"net/http"
	"nh-be/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.GetAllUsers(c.Request.Context())
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get all users", err.Error())
			return
		}
		userResp := make([]UserResponseDto, len(users))
		for i, user := range users {
			userResp[i] = UserResponseDto{
				ID: user.ID.String(),
				Username: user.Username,
				Email: user.Email,
				Role: user.Role,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}
		}
		utils.MakeSuccessResponse(c, "Users fetched successfully", users)
	}
}

func GetUserByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}
		user, err := s.GetUserById(c.Request.Context(), userID)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get user", err.Error())
			return
		}
		userResp := UserResponseDto{
			ID: user.ID.String(),
			Username: user.Username,
			Email: user.Email,
			Role: user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		utils.MakeSuccessResponse(c, "User fetched successfully", userResp)
	}
}

func GetMeHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userId = sessions.Default(c).Get("user_id")
		if userId == nil {
			utils.MakeErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User not found")
			return
		}
		
		userIdStr, ok := userId.(string)
		if !ok {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID format", "")
			return
		}
		
		parsedId, err := uuid.Parse(userIdStr)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}
		
		user, err := s.GetUserById(c.Request.Context(), parsedId)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get user", err.Error())
			return
		}
		userResp := UserResponseDto{
			ID: user.ID.String(),
			Username: user.Username,
			Email: user.Email,
			Role: user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		utils.MakeSuccessResponse(c, "User fetched successfully", userResp)
	}
}