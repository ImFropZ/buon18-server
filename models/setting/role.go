package setting

import (
	"strings"

	"system.buon18.com/m/models"

	"github.com/nullism/bqb"
)

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

func (request SettingRoleUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "description":
		bqbQuery.Comma("description = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}
	return nil
}
