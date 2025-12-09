package permission

import (
	"net/http"
	"nh-be/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Permission Handlers

func GetAllPermissionsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		permissions, count, err := s.GetAllPermissions(c.Request.Context(), search)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permissions", err.Error())
			return
		}

		var resp []PermissionResponseDto
		for _, p := range permissions {
			resp = append(resp, PermissionResponseDto{
				ID:          p.ID,
				Name:        p.Name,
				Description: p.Description,
			})
		}

		utils.MakeSuccessResponse(c, "Permissions fetched successfully", resp, count)
	}
}

func GetPermissionHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID", err.Error())
			return
		}

		p, err := s.GetPermissionByID(c.Request.Context(), id)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusNotFound, "Permission not found", err.Error())
			return
		}

		resp := PermissionResponseDto{
			ID:          p.ID.String(),
			Name:        p.Name,
			Description: p.Description,
		}

		utils.MakeSuccessResponse(c, "Permission fetched successfully", resp)
	}
}

// Permission Group Handlers
func GetAllPermissionGroupsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
		name := c.Query("name")
		assignedUser, err := uuid.Parse(c.Query("assignedUser"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid assigned user ID", err.Error())
			return
		}

		groups, count, err := s.GetAllPermissionGroups(c.Request.Context(), name, assignedUser, page, pageSize)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permission groups", err.Error())
			return
		}

		var resp []PermissionGroupResponseDto
		for _, g := range groups {
			var perms []PermissionResponseDto
			for _, p := range g.Permissions {
				perms = append(perms, PermissionResponseDto{
					ID:          p.ID.String(),
					Name:        p.Name,
					Description: p.Description,
				})
			}
			resp = append(resp, PermissionGroupResponseDto{
				ID:          g.ID.String(),
				Name:        g.Name,
				Description: g.Description,
				Permissions: perms,
				CreatedAt:   g.CreatedAt,
				UpdatedAt:   g.UpdatedAt,
			})
		}

		utils.MakeSuccessResponse(c, "Permission groups fetched successfully", resp, count)
	}
}

func CreatePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreatePermissionGroupDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		if err := s.CreatePermissionGroup(c.Request.Context(), &dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create permission group", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "Permission group created successfully", nil)
	}
}

func GetPermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID", err.Error())
			return
		}

		g, err := s.GetPermissionGroupByID(c.Request.Context(), id)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", err.Error())
			return
		}

		var perms []PermissionResponseDto
		for _, p := range g.Permissions {
			perms = append(perms, PermissionResponseDto{
				ID:          p.ID.String(),
				Name:        p.Name,
				Description: p.Description,
			})
		}
		resp := PermissionGroupResponseDto{
			ID:          g.ID.String(),
			Name:        g.Name,
			Description: g.Description,
			Permissions: perms,
			CreatedAt:   g.CreatedAt,
			UpdatedAt:   g.UpdatedAt,
		}

		utils.MakeSuccessResponse(c, "Permission group fetched successfully", resp)
	}
}

func UpdatePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID", err.Error())
			return
		}

		var dto UpdatePermissionGroupDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		if err := s.UpdatePermissionGroup(c.Request.Context(), id, &dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update permission group", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "Permission group updated successfully", nil)
	}
}

func DeletePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID", err.Error())
			return
		}

		if err := s.DeletePermissionGroup(c.Request.Context(), id); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to delete permission group", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "Permission group deleted successfully", nil)
	}
}

// User Permission Handlers

func AssignUserToGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto AssignUserGroupDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		if err := s.AssignUserToGroup(c.Request.Context(), &dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to assign user to group", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "User assigned to group successfully", nil)
	}
}

func RemoveUserFromGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("userId"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid User ID", err.Error())
			return
		}
		groupID, err := uuid.Parse(c.Param("groupId"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid Group ID", err.Error())
			return
		}

		if err := s.RemoveUserFromGroup(c.Request.Context(), userID, groupID); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to remove user from group", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "User removed from group successfully", nil)
	}
}
