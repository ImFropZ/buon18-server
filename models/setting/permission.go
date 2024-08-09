package setting

import "server/models"

type SettingPermission struct {
	*models.CommonModel
	Id   uint
	Name string
}

type SettingPermissionResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

func SettingPermissionsToResponse(permissions []SettingPermission) []SettingPermissionResponse {
	result := make([]SettingPermissionResponse, 0)
	for _, permission := range permissions {
		result = append(result, SettingPermissionResponse{
			Id:   permission.Id,
			Name: permission.Name,
		})
	}
	return result
}
