package models

type SettingPermission struct {
	*CommonModel
	Id   uint
	Name string
}

type SettingRole struct {
	*CommonModel
	Id          uint
	Name        string
	Description string
}

type SettingUser struct {
	*CommonModel
	Id    uint
	Name  string
	Email string
	Pwd   string
	Typ   string
	// -- Foreign keys
	SettingRoleId uint
}

type SettingPermissionResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type SettingRoleResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions []SettingPermissionResponse
}

type SettingUserResponse struct {
	Id    uint                `json:"id"`
	Name  string              `json:"name"`
	Email string              `json:"email"`
	Type  string              `json:"type"`
	Role  SettingRoleResponse `json:"role"`
}

func SettingPermissionsToResponse(permissions []SettingPermission) []SettingPermissionResponse {
	var result []SettingPermissionResponse
	for _, permission := range permissions {
		result = append(result, SettingPermissionResponse{
			Id:   permission.Id,
			Name: permission.Name,
		})
	}
	return result
}

func SettingUserToResponse(user SettingUser, role SettingRole, permissions []SettingPermission) SettingUserResponse {
	return SettingUserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Type:  user.Typ,
		Role: SettingRoleResponse{
			Id:          role.Id,
			Name:        role.Name,
			Description: role.Description,
			Permissions: SettingPermissionsToResponse(permissions),
		},
	}
}
