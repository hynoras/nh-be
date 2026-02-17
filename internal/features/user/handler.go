package user

import (
	"net/http"
	"nh-be/internal/constant"
	"nh-be/internal/infra"
	"nh-be/internal/utils/httputil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateUserHandler godoc
// @Summary Create a new user
// @Description Create a new user with username, email, password, and optional permissions
// @Tags Users
// @Accept json
// @Produce json
// @Param request body CreateUserDto true "User creation details"
// @Success 201 {object} httputil.SuccessResponse "User created successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid username or invalid permissions"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 409 {object} httputil.ErrorResponse "Username or email already exists"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to create user"
// @Security SessionAuth
// @Router /users [post]
func CreateUserHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreateUserDto
		var validationErr error
		validationErr = httputil.ValidateRequestFormat(c, &dto)
		if validationErr != nil {
			return
		}

		// Normalize username to lowercase and validate
		dto.Username = strings.ToLower(dto.Username)
		if err := ValidateUsername(dto.Username); err != nil {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid username", err.Error())
			return
		}

		var parsedPermissions []uuid.UUID
		if dto.Permissions != nil {
			parsedPermissions, validationErr = httputil.ParseStringsToUUIDs(dto.Permissions)
			if validationErr != nil {
				httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid permissions", validationErr.Error())
				return
			}
		}

		cleanInput := UserInput{
			Username:    dto.Username,
			Email:       dto.Email,
			Password:    dto.Password,
			Permissions: parsedPermissions,
		}
		serviceErr := s.CreateUser(c.Request.Context(), &cleanInput)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrCreateUserFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusCreated, "User created successfully", nil)
	}
}

// GetAllUsersHandler godoc
// @Summary Get all users
// @Description Retrieve a paginated list of users with optional search filter
// @Tags Users
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter users by username or email"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} httputil.SuccessResponse{data=[]UserResponseDto} "Users fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid pagination parameters"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get all users"
// @Security SessionAuth
// @Router /users [get]
func GetAllUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")

		pageInt, pageSizeInt, err := httputil.ParsePaginationParams(c)
		if err != nil {
			return
		}

		users, length, serviceErr := s.GetAllUsers(c.Request.Context(), search, pageInt, pageSizeInt)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrGetAllUsersFailed) {
			return
		}
		userResp := MapUsersToDto(users)
		httputil.MakeSuccessResponse(c, http.StatusOK, "Users fetched successfully", userResp, length)
	}
}

// GetUserByIDHandler godoc
// @Summary Get user by ID
// @Description Retrieve a single user by their UUID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID format)"
// @Success 200 {object} httputil.SuccessResponse{data=UserResponseDto} "User fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "User not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get user"
// @Security SessionAuth
// @Router /users/{id} [get]
func GetUserByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}
		user, _, serviceErr := s.GetUserById(c.Request.Context(), *parsedId, false)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrGetUserDetailFailed) {
			return
		}

		httputil.MakeSuccessResponse(c, http.StatusOK, "User fetched successfully", MapUserToDto(*user))

	}
}

// GetMeHandler godoc
// @Summary Get current user
// @Description Retrieve the currently authenticated user's information with their permissions
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} httputil.SuccessResponse{data=MeResponseDto} "User fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 401 {object} httputil.ErrorResponse "Unauthorized"
// @Failure 404 {object} httputil.ErrorResponse "User not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to get me"
// @Security SessionAuth
// @Router /users/me [get]
func GetMeHandler(s Service, sessionStore infra.SessionStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("auth_session")
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User not found")
			return
		}
		userId, err := sessionStore.GetUserSession(c.Request.Context(), cookie.Value)
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User not found")
			return
		}

		parsedUuid, err := httputil.ParseStringToUUID(userId)
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusUnauthorized, "Unauthorized", "User not found")
			return
		}

		user, perm, serviceErr := s.GetUserById(c.Request.Context(), parsedUuid, true)
		switch serviceErr {
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "User not found", serviceErr.Error())
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusOK, "User fetched successfully", MapUserToMeDto(*user, perm))
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get me", serviceErr.Error())
		}
	}
}

// UpdateUserHandler godoc
// @Summary Update a user
// @Description Update an existing user by their UUID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID format)"
// @Param request body UpdateUserDto true "Updated user details"
// @Success 200 {object} httputil.SuccessResponse "User updated successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid user ID, invalid username, or invalid permissions"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to update user"
// @Security SessionAuth
// @Router /users/{id} [put]
func UpdateUserHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := httputil.ParseStringToUUID(c.Param("id"))
		if err != nil {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
			return
		}

		var dto UpdateUserDto
		var validationErr error
		validationErr = httputil.ValidateRequestFormat(c, &dto)
		if validationErr != nil {
			return
		}

		// Validate username if being updated
		if dto.Username != "" {
			// Normalize username to lowercase and validate
			dto.Username = strings.ToLower(dto.Username)
			if err := ValidateUsername(dto.Username); err != nil {
				httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid username", err.Error())
				return
			}
		}

		var parsedPermissions []uuid.UUID
		parsedPermissions, validationErr = httputil.ParseStringsToUUIDs(dto.Permissions)
		if validationErr != nil {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid permissions", validationErr.Error())
			return
		}

		cleanInput := UserInput{
			Username:    dto.Username,
			Email:       dto.Email,
			Permissions: parsedPermissions,
		}

		serviceErr := s.UpdateUser(c.Request.Context(), userID, &cleanInput)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrUpdateUserFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "User updated successfully", nil)
	}
}

// DeleteUsersHandler godoc
// @Summary Delete users
// @Description Delete multiple users by their UUIDs
// @Tags Users
// @Accept json
// @Produce json
// @Param request body DeleteUsersDto true "Array of user IDs to delete"
// @Success 200 {object} httputil.SuccessResponse "Users deleted successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request body or invalid user ID"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to delete users"
// @Security SessionAuth
// @Router /users [delete]
func DeleteUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req DeleteUsersDto
		if err := c.ShouldBindJSON(&req); err != nil {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}
		ids := make([]uuid.UUID, len(req.IDs))
		for i, id := range req.IDs {
			parsedId, err := httputil.ParseStringToUUID(id)
			if err != nil {
				httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
				return
			}
			ids[i] = parsedId
		}
		serviceErr := s.DeleteUsers(c.Request.Context(), ids)
		if httputil.MakeServiceErrorResponse(c, serviceErr, constant.ErrDeleteUserFailed) {
			return
		}
		httputil.MakeSuccessResponse(c, http.StatusOK, "Users deleted successfully", nil)
	}
}
