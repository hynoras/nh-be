package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *interface{} `json:"data,omitempty"`
	Error   *interface{} `json:"error,omitempty"`
	Length  *int64       `json:"length,omitempty"`
}

func MakeSuccessResponse(c *gin.Context, message string, args ...interface{}) {
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
	
	c.JSON(http.StatusOK, resp)
}

func MakeErrorResponse(c *gin.Context, statusCode int, message string, error interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Error:   &error,
	})
}

func ValidateUUID(c *gin.Context, id string) (*uuid.UUID, error){
	parsedID, err := ParseStringToUUID(id)
	if err != nil {
		MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID format", err.Error())
		return nil, err
	}
	return &parsedID, nil
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