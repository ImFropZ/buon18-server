package setting

import (
	"system.buon18.com/m/models"
)

type SettingRole struct {
	*models.CommonModel
	Id          *uint
	Name        *string
	Description *string
}

func (SettingRole) AllowFilterFieldsAndOps() []string {
	return []string{"id:eq", "name:like", "name:ilike", "description:like", "description:ilike"}
}

func (SettingRole) AllowSorts() []string {
	return []string{"name"}
}

type SettingRoleResponse struct {
	Id          *uint                        `json:"id,omitempty"`
	Name        *string                      `json:"name,omitempty"`
	Description *string                      `json:"description,omitempty"`
	Permissions *[]SettingPermissionResponse `json:"permissions,omitempty"`
}

func SettingRoleToResponse(role SettingRole, permissions *[]SettingPermissionResponse) SettingRoleResponse {
	return SettingRoleResponse{
		Id:          role.Id,
		Name:        role.Name,
		Description: role.Description,
		Permissions: permissions,
	}
}

type SettingRoleCreateRequest struct {
	Name          string `json:"name" validate:"required,max=64"`
	Description   string `json:"description" validate:"required,max=255"`
	PermissionIds []uint `json:"permission_ids" validate:"required,gt=0,dive"`
}

type SettingRoleUpdateRequest struct {
	Name                *string `json:"name" validate:"omitempty,max=64"`
	Description         *string `json:"description" validate:"omitempty,max=255"`
	AddPermissionIds    *[]uint `json:"add_permission_ids" validate:"omitempty,gt=0,dive"`
	RemovePermissionIds *[]uint `json:"remove_permission_ids" validate:"omitempty,gt=0,dive"`
}
