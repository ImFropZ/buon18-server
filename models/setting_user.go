package models

var SettingUserAllowFilterFieldsAndOps = []string{"name-like", "email-like", "type-eq", "role_id-eq"}
var SettingUserAllowSortFields = []string{"name", "email", "type"}

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

type SettingUserResponse struct {
	Id    uint                `json:"id"`
	Name  string              `json:"name"`
	Email string              `json:"email"`
	Type  string              `json:"type"`
	Role  SettingRoleResponse `json:"role"`
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
