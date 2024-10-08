package setting

import (
	"system.buon18.com/m/models"
)

var SettingPermissionAllowFilterFieldsAndOps = []string{"name:like"}
var SettingPermissionAllowSortFields = []string{"name"}

type SettingPermission struct {
	*models.CommonModel
	Id   uint
	Name string
}

func (SettingPermission) AllowFilterFieldsAndOps() []string {
	return []string{"name:like"}
}

func (SettingPermission) AllowSorts() []string {
	return []string{"name"}
}

type SettingPermissionResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func SettingPermissionToResponse(permission SettingPermission) SettingPermissionResponse {
	return SettingPermissionResponse{
		Id:   permission.Id,
		Name: permission.Name,
	}
}
