package setting

import "server/models"

var SettingRoleAllowFilterFieldsAndOps = []string{"name:like", "description:like"}
var SettingRoleAllowSortFields = []string{"name"}

type SettingRole struct {
	*models.CommonModel
	Id          uint
	Name        string
	Description string
}

type SettingRoleResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Permissions []SettingPermissionResponse
}

func SettingRoleToResponse(role SettingRole, permissions []SettingPermissionResponse) SettingRoleResponse {
	return SettingRoleResponse{
		Id:          role.Id,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
	}
}
