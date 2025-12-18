package permission

func MapPermissionToDto(p Permission) PermissionResponseDto {
	return PermissionResponseDto{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
	}
}

func MapPermissionsToDto(permissions []Permission) []PermissionResponseDto {
	var result []PermissionResponseDto
	for _, p := range permissions {
		result = append(result, MapPermissionToDto(p))
	}
	return result
}

func MapPermissionGroupToDto(g PermissionGroup) PermissionGroupResponseDto {
	return PermissionGroupResponseDto{
		ID:            g.ID.String(),
		Name:          g.Name,
		Description:   g.Description,
		Permissions:   MapPermissionsToDto(g.Permissions),
		CreatedAt:     g.CreatedAt,
		UpdatedAt:     g.UpdatedAt,
	}
}

func MapPermissionGroupsToDto(groups []PermissionGroup) []PermissionGroupResponseDto {
	var result []PermissionGroupResponseDto
	for _, g := range groups {
		result = append(result, MapPermissionGroupToDto(g))
	}
	return result
}
