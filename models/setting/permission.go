package setting

import "server/models"

var SettingPermissionAllowFilterFieldsAndOps = []string{"name:like"}
var SettingPermissionAllowSortFields = []string{"name"}

type SettingPermission struct {
	*models.CommonModel
	Id   uint
	Name string
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
