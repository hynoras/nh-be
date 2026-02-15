package user

import "nh-be/internal/features/permission"

func MapUsersToDto(users []User) []UserResponseDto {
	var userDtos []UserResponseDto
	for _, user := range users {
		userDtos = append(userDtos, MapUserToListDto(user))
	}
	return userDtos
}

func MapUserToListDto(user User) UserResponseDto {
	var permissionGroups []permission.PermissionGroupResponseDto
	for _, permissionGroup := range user.AssignedPermissionGroups {
		permissionGroups = append(permissionGroups, permission.PermissionGroupResponseDto{
			ID:          permissionGroup.ID.String(),
			Name:        permissionGroup.Name,
			Description: permissionGroup.Description,
			Permissions: []permission.PermissionResponseDto{}, // Omit permissions in list view
		})
	}

	return UserResponseDto{
		ID:               user.ID.String(),
		Username:         user.Username,
		Email:            user.Email,
		PermissionGroups: permissionGroups,
		CreatedAt:        user.CreatedAt,
	}
}

func MapUserToMeDto(user User, permissionCodes []string) MeResponseDto {
	return MeResponseDto{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		Permissions: permissionCodes,
		CreatedAt:   user.CreatedAt,
	}
}

func MapUserToDto(user User) UserResponseDto {
	var permissionGroups []permission.PermissionGroupResponseDto
	for _, permissionGroup := range user.AssignedPermissionGroups {
		var permissions []permission.PermissionResponseDto
		for _, perm := range permissionGroup.Permissions {
			permissions = append(permissions, permission.PermissionResponseDto{
				ID:          perm.ID.String(),
				Name:        perm.Name,
				Description: perm.Description,
			})
		}
		permissionGroups = append(permissionGroups, permission.PermissionGroupResponseDto{
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
		PermissionGroups: permissionGroups,
		CreatedAt:        user.CreatedAt,
	}
}

func MapUserToCreatedUser(user User) CreatedUserDto {
	return CreatedUserDto{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}
}
