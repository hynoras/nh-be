package user

import (
	"net/http"
	"nh-be/constant"
	"nh-be/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUserHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreateUserDto
		var validationErr error
		validationErr = utils.ValidateRequestFormat(c, &dto)
		if validationErr != nil {
			return
		}

		var parsedPermissions []uuid.UUID
		if dto.Permissions != nil {
			parsedPermissions, validationErr = utils.ParseStringsToUUIDs(dto.Permissions)
			if validationErr != nil {
				utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid permissions", validationErr.Error())
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
		switch serviceErr {
		case ErrForbidCreateUser:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		//using the same message to avoid attacker from knowing the exact error
		case ErrDuplicateUsername:
		case ErrDuplicateEmail:
			utils.MakeErrorResponse(c, http.StatusConflict, "Username or email already exists", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusCreated, "User created successfully", nil)
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create user", serviceErr.Error())
			return
		}
	}
}

func GetAllUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")

		pageInt, pageSizeInt, err := utils.ParsePaginationParams(c)
		if err != nil {
			return
		}

		users, length, serviceErr := s.GetAllUsers(c.Request.Context(), search, pageInt, pageSizeInt)
		switch serviceErr {
		case ErrForbidViewUsers:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case nil:
			userResp := MapUsersToDto(users)
			utils.MakeSuccessResponse(c, http.StatusOK, "Users fetched successfully", userResp, length)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get all users", serviceErr.Error())
			return
		}
	}
}

func GetUserByIDHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := utils.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}
		user, serviceErr := s.GetUserById(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrForbidViewUser:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "User not found", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "User fetched successfully", MapUserToDto(*user))
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get user", serviceErr.Error())
		}
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
		switch serviceErr {
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "User not found", serviceErr.Error())
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "User fetched successfully", MapUserToDto(*user))
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get me", serviceErr.Error())
		}
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
		var validationErr error
		validationErr = utils.ValidateRequestFormat(c, &dto)
		if validationErr != nil {
			return
		}

		var parsedPermissions []uuid.UUID
		parsedPermissions, validationErr = utils.ParseStringsToUUIDs(dto.Permissions)
		if validationErr != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid permissions", validationErr.Error())
			return
		}

		cleanInput := UserInput{
			Username:    dto.Username,
			Email:       dto.Email,
			Permissions: parsedPermissions,
		}

		serviceErr := s.UpdateUser(c.Request.Context(), userID, &cleanInput)
		switch serviceErr {
		case ErrForbidUpdateUser:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "User updated successfully", nil)
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update user", serviceErr.Error())
			return
		}
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
		serviceErr := s.DeleteUsers(c.Request.Context(), ids)
		switch serviceErr {
		case ErrForbidDeleteUser:
			utils.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "Users deleted successfully", nil)
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to delete users", serviceErr.Error())
			return
		}
	}
}
