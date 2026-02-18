package httputil

import (
	"errors"
	"net/http"
	"nh-be/internal/constant"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type APIResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Error   *interface{} `json:"error,omitempty"`
	Length  *int64       `json:"length,omitempty"`
}

// for Swagger documentation
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
	Length  *int64      `json:"length,omitempty"`
}

// for Swagger documentation
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Operation failed"`
	Error   string `json:"error,omitempty" example:"Error details"`
}

type ValidationError struct {
	Index int    `json:"index"`
	Value string `json:"value"`
	Error string `json:"error"`
}

func MakeSuccessResponse(c *gin.Context, statusCode int, message string, args ...interface{}) {
	resp := APIResponse{
		Success: true,
		Message: message,
	}

	// First optional arg is data
	if len(args) > 0 && args[0] != nil {
		data := args[0]
		resp.Data = &data
	}

	// Second optional arg is length
	if len(args) > 1 {
		if length, ok := args[1].(int64); ok {
			resp.Length = &length
		} else if length, ok := args[1].(int); ok {
			length64 := int64(length)
			resp.Length = &length64
		}
	}

	c.JSON(statusCode, resp)
}

func MakeErrorResponse(c *gin.Context, statusCode int, message string, error interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   &error,
	})
}

func ValidateUUID(c *gin.Context, id string) (*uuid.UUID, error) {
	parsedID, err := ParseStringToUUID(id)
	if err != nil {
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return nil, err
	}
	return &parsedID, nil
}

func ValidateUUIDs(c *gin.Context, uuids []string) ([]uuid.UUID, error) {
	var parsedIDs []uuid.UUID
	var validationErrors []ValidationError

	for i, uuidStr := range uuids {
		// Check for empty string
		if uuidStr == "" {
			validationErrors = append(validationErrors, ValidationError{
				Index: i,
				Value: uuidStr,
				Error: "empty string is not allowed",
			})
			continue
		}

		// Attempt to parse UUID
		parsedID, err := ParseStringToUUID(uuidStr)
		if err != nil {
			validationErrors = append(validationErrors, ValidationError{
				Index: i,
				Value: uuidStr,
				Error: err.Error(),
			})
			continue
		}
		parsedIDs = append(parsedIDs, parsedID)

	}

	// If any validation errors occurred, return them
	if len(validationErrors) > 0 {
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID format", validationErrors)
		return nil, errors.New("validation failed")
	}

	return parsedIDs, nil
}

func ValidateRequestFormat(c *gin.Context, dto interface{}) error {
	if err := c.ShouldBindJSON(&dto); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			MakeErrorResponse(
				c,
				http.StatusUnprocessableEntity,
				"Validation failed",
				verr.Error(),
			)
			return verr
		}
		MakeErrorResponse(
			c,
			http.StatusBadRequest,
			"Invalid request",
			err.Error(),
		)
		return err
	}
	return nil
}

func MakeServiceErrorResponse(c *gin.Context, err error, msg string) bool {
	if err == nil {
		return false
	}

	switch err {
	//auth
	case constant.ErrInvalidCredentials:
		MakeErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", err.Error())
	case constant.ErrVerificationTokenNotFound:
		MakeErrorResponse(c, http.StatusNotFound, "Verification token not found", err.Error())
	case constant.ErrUnauthenticated:
		MakeErrorResponse(c, http.StatusUnauthorized, "Unauthenticated", err.Error())
	case constant.ErrEmailAlreadyExists:
		MakeErrorResponse(c, http.StatusConflict, "Email already exists", err.Error())
	case constant.ErrVerificationTokenExpired:
		MakeErrorResponse(c, http.StatusUnauthorized, "Verification token expired", err.Error())
	case constant.ErrInvalidVerificationToken:
		MakeErrorResponse(c, http.StatusUnauthorized, "Invalid verification token", err.Error())
	case constant.ErrSessionNotFound:
		MakeErrorResponse(c, http.StatusUnauthorized, "Session not found", err.Error())

	//permission
	case constant.ErrPermissionNotFound:
		MakeErrorResponse(c, http.StatusNotFound, "Permission not found", err.Error())
	case constant.ErrPermissionGroupNotFound:
		MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", err.Error())
	case constant.ErrNotNullPermissions:
		MakeErrorResponse(c, http.StatusBadRequest, "Permissions can not be null", err.Error())
	case constant.ErrCannotDeleteSuperAdmin:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrDeletePermissionFailed, err.Error())
	case constant.ErrForbidViewPermissions:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetAllPermissionFailed, err.Error())
	case constant.ErrPermissionGroupNameAlreadyExists:
		MakeErrorResponse(c, http.StatusConflict, constant.ErrCreatePermissionGroupFailed, err.Error())
	case constant.ErrForbidViewPermissionGroups:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetAllPermissionGroupFailed, err.Error())
	case constant.ErrForbidViewPermissionGroup:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetPermissionGroupDetailFailed, err.Error())
	case constant.ErrForbidCreatePermissionGroup:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrCreatePermissionGroupFailed, err.Error())
	case constant.ErrForbidUpdatePermissionGroup:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrUpdatePermissionGroupFailed, err.Error())
	case constant.ErrForbidDeletePermissionGroup:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrDeletePermissionGroupFailed, err.Error())

	//user
	case constant.ErrUserNotFound:
		MakeErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
	case constant.ErrForbidViewUsers:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetAllUsersFailed, err.Error())
	case constant.ErrForbidViewUser:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetUserDetailFailed, err.Error())
	case constant.ErrForbidUpdateUser:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrUpdateUserFailed, err.Error())
	case constant.ErrForbidDeleteUser:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrDeleteUserFailed, err.Error())
	case constant.ErrDuplicateUsername:
		MakeErrorResponse(c, http.StatusConflict, "Invalid username", err.Error())
	case constant.ErrDuplicateEmail:
		MakeErrorResponse(c, http.StatusConflict, "Invalid email", err.Error())
	case constant.ErrUsernameMustStartWithLetter:
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid username", err.Error())
	case constant.ErrUsernameMustEndWithLetterOrNumber:
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid username", err.Error())
	case constant.ErrUsernameNoConsecutiveSpecialChars:
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid username", err.Error())
	case constant.ErrUsernameNoAdjacentSpecialChars:
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid username", err.Error())
	case constant.ErrReservedUsername:
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid username", err.Error())

	//experiment
	case constant.ErrForbidViewExperiments:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetAllExperimentFailed, err.Error())
	case constant.ErrForbidViewExperiment:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetExperimentDetailFailed, err.Error())
	case constant.ErrForbidUpdateExperiment:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrUpdateExperimentFailed, err.Error())
	case constant.ErrForbidDeleteExperiment:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, err.Error())
	case constant.ErrStatusTransitionFromDraftToPlanning:
		MakeErrorResponse(c, http.StatusBadRequest, constant.ErrInvalidStatusTransition, err.Error())
	case constant.ErrStatusTransitionFromPlanningToRunning:
		MakeErrorResponse(c, http.StatusBadRequest, constant.ErrInvalidStatusTransition, err.Error())
	case constant.ErrStatusTransitionFromRunningToCompletedOrAborted:
		MakeErrorResponse(c, http.StatusBadRequest, constant.ErrInvalidStatusTransition, err.Error())
	case constant.ErrExperimentConflict:
		MakeErrorResponse(c, http.StatusConflict, constant.ErrUpdateExperimentFailed, err.Error())
	case constant.ErrExperimentAlreadyInTargetState:
		MakeErrorResponse(c, http.StatusBadRequest, constant.ErrUpdateExperimentFailed, err.Error())
	case constant.ErrExperimentNotFound:
		MakeErrorResponse(c, http.StatusNotFound, "Experiment not found", err.Error())

	//experiment result
	case constant.ErrExperimentResultNotFound:
		MakeErrorResponse(c, http.StatusNotFound, "Experiment result not found", err.Error())
	case constant.ErrExperimentResultAlreadyExists:
		MakeErrorResponse(c, http.StatusConflict, "Experiment result already exists", err.Error())
	case constant.ErrForbidCreateExperimentResult:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrCreateExperimentFailed, err.Error())
	case constant.ErrForbidViewExperimentResult:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrGetExperimentDetailFailed, err.Error())
	case constant.ErrForbidUpdateExperimentResult:
		MakeErrorResponse(c, http.StatusForbidden, constant.ErrUpdateExperimentFailed, err.Error())
	case constant.ErrInvalidOutcome:
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid outcome value", err.Error())
	case constant.ErrInvalidConfidenceLevel:
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid confidence level value", err.Error())
	case constant.ErrExperimentResultConflict:
		MakeErrorResponse(c, http.StatusConflict, constant.ErrUpdateExperimentFailed, err.Error())
	//add other domain error here

	default:
		MakeErrorResponse(c, http.StatusInternalServerError, msg, err.Error())
	}
	return true
}
