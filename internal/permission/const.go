package permission

import "errors"

var (
	ErrPermissionNotFound               = errors.New("permission not found")
	ErrPermissionGroupNotFound          = errors.New("permission group not found")
	ErrNotNullPermissions               = errors.New("permissions can not be null")
	ErrCannotDeleteSuperAdmin           = errors.New("can not delete Super Admin. At least one must exist")
	ErrPermissionGroupNameAlreadyExists = errors.New("permission group name already exists")
	ErrForbidViewPermissionGroups       = errors.New("you are not allowed to view permission groups")
	ErrForbidViewPermissionGroup        = errors.New("you are not allowed to view permission group")
	ErrForbidCreatePermissionGroup      = errors.New("you are not allowed to create permission group")
	ErrForbidUpdatePermissionGroup      = errors.New("you are not allowed to update permission group")
	ErrForbidDeletePermissionGroup      = errors.New("you are not allowed to delete permission group")
)
