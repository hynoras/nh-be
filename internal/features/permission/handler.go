package permission

import (
	"net/http"
	"nh-be/constant"
	"nh-be/internal/httputil"
	"nh-be/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Permission Handlers

// GetAllPermissionsHandler godoc
// @Summary Get all permissions
// @Description Retrieve a list of all permissions with optional search filter
// @Tags Permissions
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter permissions by name"
// @Success 200 {object} httputil.SuccessResponse{data=[]PermissionResponseDto} "Permissions fetched successfully"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to fetch permissions"
// @Security SessionAuth
// @Router /permissions [get]
func GetAllPermissionsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		permissions, count, serviceErr := s.GetAllPermissions(c.Request.Context(), search)
		switch serviceErr {
		case ErrForbidViewPermissions:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case nil:
			resp := MapPermissionsToDto(permissions)
			httputil.MakeSuccessResponse(c, http.StatusOK, "Permissions fetched successfully", resp, count)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permissions", serviceErr.Error())
			return
		}
	}
}

// GetPermissionHandler godoc
// @Summary Get permission by ID
// @Description Retrieve a single permission by its UUID
// @Tags Permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID (UUID format)"
// @Success 200 {object} httputil.SuccessResponse{data=PermissionResponseDto} "Permission fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 404 {object} httputil.ErrorResponse "Permission not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to fetch permission"
// @Security SessionAuth
// @Router /permissions/{id} [get]
func GetPermissionHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		permission, serviceErr := s.GetPermissionByID(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrPermissionNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Permission not found", serviceErr.Error())
			return
		case nil:
			resp := MapPermissionToDto(*permission)
			httputil.MakeSuccessResponse(c, http.StatusOK, "Permission fetched successfully", resp)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permission", serviceErr.Error())
			return
		}
	}
}

// Permission Group Handlers

// CreatePermissionGroupHandler godoc
// @Summary Create a new permission group
// @Description Create a new permission group (role) with specified permissions
// @Tags Permission Groups
// @Accept json
// @Produce json
// @Param request body CreatePermissionGroupDto true "Permission group details"
// @Success 201 {object} httputil.SuccessResponse "Permission group created successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request format or invalid permissions"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 409 {object} httputil.ErrorResponse "Role name already exists"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to create permission group"
// @Security SessionAuth
// @Router /permission-groups [post]
func CreatePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto CreatePermissionGroupDto
		var validationErr error
		validationErr = httputil.ValidateRequestFormat(c, &dto)
		if validationErr != nil {
			return
		}

		var parsedPermissions []uuid.UUID
		parsedPermissions, validationErr = httputil.ValidateUUIDs(c, dto.Permissions)
		if validationErr != nil {
			return
		}

		cleanInput := PermissionGroupInput{
			Name:        dto.Name,
			Description: dto.Description,
			Permissions: parsedPermissions,
		}

		serviceErr := s.CreatePermissionGroup(c.Request.Context(), &cleanInput)
		switch serviceErr {
		case ErrForbidCreatePermissionGroup:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case ErrPermissionGroupNameAlreadyExists:
			httputil.MakeErrorResponse(c, http.StatusConflict, "Role name already exists", serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusCreated, "Permission group created successfully", nil)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to create permission group", serviceErr.Error())
			return
		}
	}
}

// GetAllPermissionGroupsHandler godoc
// @Summary Get all permission groups
// @Description Retrieve a paginated list of permission groups with optional filters
// @Tags Permission Groups
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter permission groups by name"
// @Param permissions query []string false "Filter by permission IDs (UUID format)"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} httputil.SuccessResponse{data=[]PermissionGroupResponseDto} "Permission groups fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid pagination parameters or permission IDs"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to fetch permission groups"
// @Security SessionAuth
// @Router /permission-groups [get]
func GetAllPermissionGroupsHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, pageSize, err := httputil.ParsePaginationParams(c)
		if err != nil {
			return
		}
		search := c.Query("search")

		var permissionIds []uuid.UUID
		permissionIdStrings := c.QueryArray("permissions")
		if len(permissionIdStrings) > 0 {
			parsedPermissionIds, validationErr := httputil.ValidateUUIDs(c, permissionIdStrings)
			if validationErr != nil {
				return
			}
			permissionIds = parsedPermissionIds
		}

		groups, count, serviceErr := s.GetAllPermissionGroups(c.Request.Context(), search, permissionIds, page, pageSize)
		switch serviceErr {
		case ErrForbidViewPermissionGroups:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case nil:
			resp := MapPermissionGroupsToDto(groups)
			httputil.MakeSuccessResponse(c, http.StatusOK, "Permission groups fetched successfully", resp, count)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permission groups", serviceErr.Error())
			return
		}
	}
}

// GetPermissionGroupHandler godoc
// @Summary Get permission group by ID
// @Description Retrieve a single permission group by its UUID
// @Tags Permission Groups
// @Accept json
// @Produce json
// @Param id path string true "Permission Group ID (UUID format)"
// @Success 200 {object} httputil.SuccessResponse{data=PermissionGroupResponseDto} "Permission group fetched successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Permission group not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to fetch permission group"
// @Security SessionAuth
// @Router /permission-groups/{id} [get]
func GetPermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		permissionGroup, serviceErr := s.GetPermissionGroupByID(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrForbidViewPermissionGroup:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", nil)
			return
		case nil:
			resp := MapPermissionGroupToDto(*permissionGroup)
			httputil.MakeSuccessResponse(c, http.StatusOK, "Permission group fetched successfully", resp)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to fetch permission group", serviceErr.Error())
			return
		}
	}
}

// UpdatePermissionGroupHandler godoc
// @Summary Update a permission group
// @Description Update an existing permission group by its UUID
// @Tags Permission Groups
// @Accept json
// @Produce json
// @Param id path string true "Permission Group ID (UUID format)"
// @Param request body UpdatePermissionGroupDto true "Updated permission group details"
// @Success 200 {object} httputil.SuccessResponse "Permission group updated successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format or invalid permissions"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Permission group not found"
// @Failure 409 {object} httputil.ErrorResponse "Role name already exists"
// @Failure 422 {object} httputil.ErrorResponse "Validation failed"
// @Failure 500 {object} httputil.ErrorResponse "Failed to update permission group"
// @Security SessionAuth
// @Router /permission-groups/{id} [put]
func UpdatePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		var dto UpdatePermissionGroupDto
		var validationErr error
		validationErr = httputil.ValidateRequestFormat(c, &dto)
		if validationErr != nil {
			return
		}

		var parsedPermissions []uuid.UUID
		parsedPermissions, validationErr = utils.ParseStringsToUUIDs(dto.Permissions)
		if validationErr != nil {
			httputil.MakeErrorResponse(c, http.StatusBadRequest, "Invalid permissions", validationErr.Error())
			return
		}

		cleanInput := PermissionGroupInput{
			Name:        dto.Name,
			Description: dto.Description,
			Permissions: parsedPermissions,
		}

		serviceErr := s.UpdatePermissionGroup(c.Request.Context(), *parsedId, &cleanInput)
		switch serviceErr {
		case ErrForbidUpdatePermissionGroup:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", nil)
			return
		case ErrPermissionGroupNameAlreadyExists:
			httputil.MakeErrorResponse(c, http.StatusConflict, "Role name already exists", serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusOK, "Permission group updated successfully", nil)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to update permission group", serviceErr.Error())
			return
		}
	}
}

// DeletePermissionGroupHandler godoc
// @Summary Delete a permission group
// @Description Delete an existing permission group by its UUID
// @Tags Permission Groups
// @Accept json
// @Produce json
// @Param id path string true "Permission Group ID (UUID format)"
// @Success 200 {object} httputil.SuccessResponse "Permission group deleted successfully"
// @Failure 400 {object} httputil.ErrorResponse "Invalid ID format"
// @Failure 403 {object} httputil.ErrorResponse "Authorization failed"
// @Failure 404 {object} httputil.ErrorResponse "Permission group not found"
// @Failure 500 {object} httputil.ErrorResponse "Failed to delete permission group or cannot delete Super Admin"
// @Security SessionAuth
// @Router /permission-groups/{id} [delete]
func DeletePermissionGroupHandler(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedId, idErr := httputil.ValidateUUID(c, c.Param("id"))
		if idErr != nil {
			return
		}

		serviceErr := s.DeletePermissionGroup(c.Request.Context(), *parsedId)
		switch serviceErr {
		case ErrForbidDeletePermissionGroup:
			httputil.MakeErrorResponse(c, http.StatusForbidden, constant.ErrAuthorizationFailed, serviceErr.Error())
			return
		case gorm.ErrRecordNotFound:
			httputil.MakeErrorResponse(c, http.StatusNotFound, "Permission group not found", nil)
			return
		case ErrCannotDeleteSuperAdmin:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Can not delete Super Admin. At least one must exist", serviceErr.Error())
			return
		case nil:
			httputil.MakeSuccessResponse(c, http.StatusOK, "Permission group deleted successfully", nil)
			return
		default:
			httputil.MakeErrorResponse(c, http.StatusInternalServerError, "Failed to delete permission group", serviceErr.Error())
			return
		}
	}
}
