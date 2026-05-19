package authutil

import (
	"context"
	"slices"
	"nh-be/internal/utils/ctxutil"

	"github.com/google/uuid"
)

type PermissionGetter interface {
	GetUserPermissionCodeNames(ctx context.Context, userId uuid.UUID) ([]string, error)
}

// RequirePermission verifies if the user from the context has at least one of the allowed permissions.
// Returns the specified forbidError if not allowed.
func RequirePermission(ctx context.Context, getter PermissionGetter, forbidError error, allowedPermissions ...string) error {
	userId, err := ctxutil.GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	userPerms, err := getter.GetUserPermissionCodeNames(ctx, userId)
	if err != nil {
		return err
	}

	for _, allowed := range allowedPermissions {
		if slices.Contains(userPerms, allowed) {
			return nil
		}
	}

	return forbidError
}
