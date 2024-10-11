package setting

import (
	"system.buon18.com/m/models"
)

type SettingPermission struct {
	*models.CommonModel
	Id   *uint
	Name *string
}

func (SettingPermission) AllowFilterFieldsAndOps() []string {
	return []string{"name:like", "name:ilike"}
}

func (SettingPermission) AllowSorts() []string {
	return []string{"name"}
}

type SettingPermissionResponse struct {
	Id   *uint   `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

func SettingPermissionToResponse(permission SettingPermission) SettingPermissionResponse {
	return SettingPermissionResponse{
		Id:   permission.Id,
		Name: permission.Name,
	}
}
