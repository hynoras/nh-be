package permission

import (
	"net/http"
	"nh-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

		resp := MapPermissionsToDto(permissions)
		utils.MakeSuccessResponse(c, http.StatusOK, "Permissions fetched successfully", resp, count)
	}
}

func GetPermissionHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := utils.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		permission, serviceErr := s.GetPermissionByID(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrPermissionNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "Permission not found", serviceErr.Error())
			return
		case nil:
			resp := MapPermissionToDto(*permission)
			utils.MakeSuccessResponse(c, http.StatusOK, "Permission fetched successfully", resp)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permission", serviceErr.Error())
			return
		}
	}
}

// Permission Group Handlers
func CreatePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreatePermissionGroupDto
		var validationErr error
		validationErr = utils.ValidateRequestFormat(c, &dto)
        if validationErr != nil {
            return
        }
		
		var parsedPermissions []uuid.UUID
		parsedPermissions, validationErr = utils.ValidateUUIDs(c, dto.Permissions)
        if validationErr != nil {
            return
        }

		cleanInput := PermissionGroupInput{
			Name: dto.Name,
			Description: dto.Description,
			Permissions: parsedPermissions,
		}

		serviceErr := s.CreatePermissionGroup(c.Request.Context(), &cleanInput)
		switch serviceErr {
		case ErrRoleNameAlreadyExists:
			utils.MakeErrorResponse(c, http.StatusConflict, "Role name already exists", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusCreated, "Permission group created successfully", nil)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create permission group", serviceErr.Error())
			return
		}
	}
}

func GetAllPermissionGroupsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, pageSize, err := utils.ParsePaginationParams(c)
		if err != nil {
			return
		}
		search := c.Query("search")

		var permissionIds []uuid.UUID
		permissionIdStrings := c.QueryArray("permissions")
		if len(permissionIdStrings) > 0 {
			parsedPermissionIds, validationErr := utils.ValidateUUIDs(c, permissionIdStrings)
			if validationErr != nil {
				return
			}
			permissionIds = parsedPermissionIds
		}

		groups, count, serviceErr := s.GetAllPermissionGroups(c.Request.Context(), search, permissionIds, page, pageSize)
		if serviceErr != nil {
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permission groups", serviceErr.Error())
			return
		}

		resp := MapPermissionGroupsToDto(groups)
		utils.MakeSuccessResponse(c, http.StatusOK, "Permission groups fetched successfully", resp, count)
	}
}

func GetPermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := utils.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		permissionGroup, serviceErr := s.GetPermissionGroupByID(c.Request.Context(), *parsedId)
		switch serviceErr {
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", nil)
			return
		case nil :
			resp := MapPermissionGroupToDto(*permissionGroup)
			utils.MakeSuccessResponse(c, http.StatusOK, "Permission group fetched successfully", resp)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permission group", serviceErr.Error())
			return
		}
	}
}

func UpdatePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := utils.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		var dto UpdatePermissionGroupDto
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

		cleanInput := PermissionGroupInput{
			Name: dto.Name,
			Description: dto.Description,
			Permissions: parsedPermissions,
		}

		serviceErr := s.UpdatePermissionGroup(c.Request.Context(), *parsedId, &cleanInput)
		switch serviceErr {
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", nil)
			return
		case ErrRoleNameAlreadyExists:
			utils.MakeErrorResponse(c, http.StatusConflict, "Role name already exists", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "Permission group updated successfully", nil)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update permission group", serviceErr.Error())
			return
		}
	}
}

func DeletePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := utils.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		serviceErr := s.DeletePermissionGroup(c.Request.Context(), *parsedId)
		switch serviceErr {
		case gorm.ErrRecordNotFound:
			utils.MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", nil)
			return
		case ErrCannotDeleteSuperAdmin:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Can not delete Super Admin. At least one must exist", serviceErr.Error())
			return
		case nil:
			utils.MakeSuccessResponse(c, http.StatusOK, "Permission group deleted successfully", nil)
			return
		default:
			utils.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to delete permission group", serviceErr.Error())
			return
		}
	}
}
