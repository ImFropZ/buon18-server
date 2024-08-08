package models

type SettingRole struct {
	*CommonModel
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
