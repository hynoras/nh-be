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
		search := c.Query("search")
		users, err := s.GetAllUsers(c.Request.Context(), search)
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
		userID, err := utils.ParseStringToUUID(c.Param("id"))
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

func UpdateUserHander(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.ParseStringToUUID(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}

		var req UpdateUserDto
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		
		userReq := &User{
			Username: req.Username,
			Email: req.Email,
			Role: req.Role,
		}

		err = s.UpdateUser(c.Request.Context(), userID, userReq)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err.Error())
			return
		}
		utils.MakeSuccessResponse(c, "User updated successfully", nil)
	}
}