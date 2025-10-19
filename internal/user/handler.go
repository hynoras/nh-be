package user

import (
	"net/http"
	"nh-be/utils"

	"github.com/gin-gonic/gin"
)

func GetAllUsersHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := s.GetAllUsers(c.Request.Context())
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to get all users", err.Error())
			return
		}
		utils.MakeSuccessResponse(c, "Users fetched successfully", users)
	}
}