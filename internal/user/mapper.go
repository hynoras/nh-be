package user

func MapUsersToDto(users []User) []UserResponseDto {
	var userDtos []UserResponseDto
	for _, user := range users {
		userDtos = append(userDtos, MapUserToListDto(user))
	}
	return userDtos
}

func MapUserToListDto(user User) UserResponseDto {
	var permissionGroups []PermissionGroupResponseDto
	for _, permissionGroup := range user.AssignedPermissionGroups {
		permissionGroups = append(permissionGroups, PermissionGroupResponseDto{
			ID:          permissionGroup.ID.String(),
			Name:        permissionGroup.Name,
			Description: permissionGroup.Description,
			Permissions: []PermissionResponseDto{}, // Omit permissions in list view
		})
	}

	return UserResponseDto{
		ID:               user.ID.String(),
		Username:         user.Username,
		Email:            user.Email,
		Role:             user.Role,
		PermissionGroups: permissionGroups,
		CreatedAt:        user.CreatedAt,
	}
}

func MapUserToDto(user User) UserResponseDto {
	var permissionGroups []PermissionGroupResponseDto
	for _, permissionGroup := range user.AssignedPermissionGroups {
		var permissions []PermissionResponseDto
		for _, permission := range permissionGroup.Permissions {
			permissions = append(permissions, PermissionResponseDto{
				ID:          permission.ID.String(),
				Name:        permission.Name,
				Description: permission.Description,
			})
		}
		permissionGroups = append(permissionGroups, PermissionGroupResponseDto{
			ID:          permissionGroup.ID.String(),
			Name:        permissionGroup.Name,
			Description: permissionGroup.Description,
			Permissions: permissions,
		})
	}

	return UserResponseDto{
		ID:               user.ID.String(),
		Username:         user.Username,
		Email:            user.Email,
		Role:             user.Role,
		PermissionGroups: permissionGroups,
		CreatedAt:        user.CreatedAt,
	}
}
