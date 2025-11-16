package user

import (
	"net/http"
	"nh-be/utils"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		role := c.Query("role")

		page := c.Query("page")
		pageSize := c.Query("page_size")
		
		var pageInt int
		var err error
		if page == "" {
			pageInt = 1
		} else {
			pageInt, err = strconv.Atoi(page)
			if err != nil {
				utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid page", err.Error())
				return
			}
		}

		var pageSizeInt int
		if pageSize == "" {
			pageSizeInt = 10
		} else {
			pageSizeInt, err = strconv.Atoi(pageSize)
			if err != nil {
				utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid page size", err.Error())
				return
			}
		}
		
		users, length, err := s.GetAllUsers(c.Request.Context(), search, role, pageInt, pageSizeInt)
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
		utils.MakeSuccessResponse(c, "Users fetched successfully", userResp, length)
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

func CreateUserHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserDto
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}
		if req.Password != req.ConfirmPassword {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Password and confirm password do not match", "")
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to hash password", err.Error())
			return
		}
		
		userReq := &User{
			Username: req.Username,
			Email: req.Email,
			Password: string(hashedPassword),
			Role: req.Role,
		}

		err = s.CreateUser(c.Request.Context(), userReq)
		if err != nil {
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