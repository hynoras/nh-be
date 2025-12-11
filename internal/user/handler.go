package user

import (
	"errors"
	"net/http"
	"nh-be/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetAllUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		role := c.Query("role")

		pageInt, pageSizeInt, err := utils.ParsePaginationParams(c, 1, 10)
		if err != nil {
			return
		}
		
		users, length, serviceErr := s.GetAllUsers(c.Request.Context(), search, role, pageInt, pageSizeInt)
		if serviceErr != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get all users", serviceErr.Error())
			return
		}
		userResp := MapUsersToDto(users)
		utils.MakeSuccessResponse(c, "Users fetched successfully", userResp, length)
	}
}

func GetUserByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := utils.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}
		user, serviceErr := s.GetUserById(c.Request.Context(), *parsedId)
		if serviceErr != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get user", serviceErr.Error())
			return
		}
		utils.MakeSuccessResponse(c, "User fetched successfully", MapUserToDto(*user))
	}
}

func GetMeHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userId = sessions.Default(c).Get("user_id")
		if userId == nil {
			utils.MakeErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User not found")
			return
		}
		
		parsedId, idErr := utils.ValidateUUID(c, userId.(string))
		if idErr != nil {
			return
		}
		
		user, serviceErr := s.GetUserById(c.Request.Context(), *parsedId)
		if serviceErr != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get user", serviceErr.Error())
			return
		}
		utils.MakeSuccessResponse(c, "User fetched successfully", MapUserToDto(*user))
	}
}

func CreateUserHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreateUserDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		err := s.CreateUser(c.Request.Context(), &dto)
		if err != nil {
			if errors.Is(err, ErrDuplicateUsername) || errors.Is(err, ErrDuplicateEmail) {
				utils.MakeErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
				return
			}
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err.Error())
			return
		}
		utils.MakeSuccessResponse(c, "User created successfully", nil)
	}
}

func UpdateUserHander(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := utils.ParseStringToUUID(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}

		var dto UpdateUserDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}


		err = s.UpdateUser(c.Request.Context(), userID, &dto)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err.Error())
			return
		}
		utils.MakeSuccessResponse(c, "User updated successfully", nil)
	}
}

func DeleteUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteUsersDto
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}
		ids := make([]uuid.UUID, len(req.IDs))
		for i, id := range req.IDs {
			parsedId, err := utils.ParseStringToUUID(id)
			if err != nil {
				utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
				return
			}
			ids[i] = parsedId
		}
		err := s.DeleteUsers(c.Request.Context(), ids)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to delete users", err.Error())
			return
		}
		utils.MakeSuccessResponse(c, "Users deleted successfully", nil)
	}
}