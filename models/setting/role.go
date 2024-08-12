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
	Id          uint                        `json:"id"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Permissions []SettingPermissionResponse `json:"permissions"`
}

func SettingRoleToResponse(role SettingRole, permissions []SettingPermissionResponse) SettingRoleResponse {
	return SettingRoleResponse{
		Id:          role.Id,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
	}
}

type SettingRoleCreateRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description" validate:"required,max=255"`
	PermissionIds []uint `json:"permission_ids" validate:"required,gt=0,dive"`
}
