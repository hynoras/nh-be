package httputil

import (
	"errors"
	"log/slog"
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
	if err := c.ShouldBindJSON(dto); err != nil {
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

// HTTPErrorMapping maps a domain error to an HTTP status code and user-facing message.
type HTTPErrorMapping struct {
	Error      error
	StatusCode int
	Message    string
}

// registry holds all registered error-to-HTTP mappings.
// Populated during init() and read-only at runtime.
var registry []HTTPErrorMapping

// RegisterError adds a domain error mapping to the global registry.
// Must only be called during package initialization (init blocks).
func RegisterError(err error, statusCode int, message string) {
	registry = append(registry, HTTPErrorMapping{
		Error:      err,
		StatusCode: statusCode,
		Message:    message,
	})
}

func MakeServiceErrorResponse(c *gin.Context, err error, msg string) bool {
	if err == nil {
		return false
	}

	for _, mapping := range registry {
		if errors.Is(err, mapping.Error) {
			MakeErrorResponse(c, mapping.StatusCode, mapping.Message, err.Error())
			return true
		}
	}

	slog.Error("unhandled service error", "error", err.Error(),
		"request_id", c.GetString(string(constant.CtxRequestId)))
	MakeErrorResponse(c, http.StatusInternalServerError, msg, "An internal error occurred")
	return true
}
