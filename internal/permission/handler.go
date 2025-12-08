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
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

		permissions, count, err := s.GetAllPermissions(c.Request.Context(), page, pageSize)
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permissions", err.Error())
			return
		}

		var resp []PermissionResponseDto
		for _, p := range permissions {
			resp = append(resp, PermissionResponseDto{
				ID:          p.ID.String(),
				Name:        p.Name,
				Description: p.Description,
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
			})
		}

		utils.MakeSuccessResponse(c, "Permissions fetched successfully", resp, count)
	}
}

func CreatePermissionHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreatePermissionDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		if err := s.CreatePermission(c.Request.Context(), &dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create permission", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "Permission created successfully", nil)
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
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		}

		utils.MakeSuccessResponse(c, "Permission fetched successfully", resp)
	}
}

func UpdatePermissionHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID", err.Error())
			return
		}

		var dto UpdatePermissionDto
		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid request", err.Error())
			return
		}

		if err := s.UpdatePermission(c.Request.Context(), id, &dto); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update permission", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "Permission updated successfully", nil)
	}
}

func DeletePermissionHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			utils.MakeErrorResponse(c, http.StatusBadRequest, "Invalid ID", err.Error())
			return
		}

		if err := s.DeletePermission(c.Request.Context(), id); err != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to delete permission", err.Error())
			return
		}

		utils.MakeSuccessResponse(c, "Permission deleted successfully", nil)
	}
}

// Permission Group Handlers

func GetAllPermissionGroupsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

		groups, count, err := s.GetAllPermissionGroups(c.Request.Context(), page, pageSize)
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
					CreatedAt:   p.CreatedAt,
					UpdatedAt:   p.UpdatedAt,
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
				CreatedAt:   p.CreatedAt,
				UpdatedAt:   p.UpdatedAt,
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
