package setting

import (
	"server/models"
	"strings"

	"github.com/nullism/bqb"
)

var SettingUserAllowFilterFieldsAndOps = []string{"name:like", "email:like", "typ:in", "role_id:eq"}
var SettingUserAllowSortFields = []string{"name", "email", "type"}

type SettingUser struct {
	*models.CommonModel
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

func SettingUserToResponse(user SettingUser, role SettingRoleResponse) SettingUserResponse {
	return SettingUserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Type:  user.Typ,
		Role:  role,
	}
}

type SettingUserCreateRequest struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	RoleId uint   `json:"role_id" validate:"required"`
}

type SettingUserUpdateRequest struct {
	Name     *string `json:"name"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password"`
	RoleId   *uint   `json:"role_id"`
}

func (request SettingUserUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldName string, value interface{}) error {

	switch strings.ToLower(fieldName) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "email":
		bqbQuery.Comma("email = ?", value)
	case "roleid":
		bqbQuery.Comma("setting_role_id = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}

	return nil
}
