package permission

import (
	"errors"
	"net/http"
	"nh-be/internal/utils/httputil"
)

const (
	ErrGetAllPermissionFailed         = "Failed to get permissions"
	ErrGetPermissionDetailFailed      = "Failed to get permission detail"
	ErrCreatePermissionFailed         = "Failed to create permission"
	ErrUpdatePermissionFailed         = "Failed to update permission"
	ErrDeletePermissionFailed         = "Failed to delete permission"
	ErrGetAllPermissionGroupFailed    = "Failed to get permission groups"
	ErrGetPermissionGroupDetailFailed = "Failed to get permission group detail"
	ErrCreatePermissionGroupFailed    = "Failed to create permission group"
	ErrUpdatePermissionGroupFailed    = "Failed to update permission group"
	ErrDeletePermissionGroupFailed    = "Failed to delete permission group"
)

var (
	ErrPermissionNotFound               = errors.New("permission not found")
	ErrPermissionGroupNotFound          = errors.New("permission group not found")
	ErrNotNullPermissions               = errors.New("permissions can not be null")
	ErrCannotDeleteSuperAdmin           = errors.New("can not delete Super Admin. At least one must exist")
	ErrForbidViewPermissions            = errors.New("you are not allowed to view permissions")
	ErrPermissionGroupNameAlreadyExists = errors.New("permission group name already exists")
	ErrForbidViewPermissionGroups       = errors.New("you are not allowed to view permission groups")
	ErrForbidViewPermissionGroup        = errors.New("you are not allowed to view permission group")
	ErrForbidCreatePermissionGroup      = errors.New("you are not allowed to create permission group")
	ErrForbidUpdatePermissionGroup      = errors.New("you are not allowed to update permission group")
	ErrForbidDeletePermissionGroup      = errors.New("you are not allowed to delete permission group")
)

func init() {
	httputil.RegisterError(ErrPermissionNotFound, http.StatusNotFound, "Permission not found")
	httputil.RegisterError(ErrPermissionGroupNotFound, http.StatusNotFound, "Permission group not found")
	httputil.RegisterError(ErrNotNullPermissions, http.StatusBadRequest, "Permissions can not be null")
	httputil.RegisterError(ErrCannotDeleteSuperAdmin, http.StatusForbidden, ErrDeletePermissionFailed)
	httputil.RegisterError(ErrForbidViewPermissions, http.StatusForbidden, ErrGetAllPermissionFailed)
	httputil.RegisterError(ErrPermissionGroupNameAlreadyExists, http.StatusConflict, ErrCreatePermissionGroupFailed)
	httputil.RegisterError(ErrForbidViewPermissionGroups, http.StatusForbidden, ErrGetAllPermissionGroupFailed)
	httputil.RegisterError(ErrForbidViewPermissionGroup, http.StatusForbidden, ErrGetPermissionGroupDetailFailed)
	httputil.RegisterError(ErrForbidCreatePermissionGroup, http.StatusForbidden, ErrCreatePermissionGroupFailed)
	httputil.RegisterError(ErrForbidUpdatePermissionGroup, http.StatusForbidden, ErrUpdatePermissionGroupFailed)
	httputil.RegisterError(ErrForbidDeletePermissionGroup, http.StatusForbidden, ErrDeletePermissionGroupFailed)
}
