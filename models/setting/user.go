package setting

import (
	"system.buon18.com/m/models"
)

type SettingUser struct {
	*models.CommonModel
	Id    *uint
	Name  *string
	Email *string
	Pwd   *string
	Typ   *string
	// -- Foreign keys
	SettingRoleId *uint
}

func (SettingUser) AllowFilterFieldsAndOps() []string {
	return []string{"name:ilike", "email:ilike", "name:like", "email:like", "typ:in", "role-id:eq"}
}

func (SettingUser) AllowSorts() []string {
	return []string{"name", "email", "typ"}
}

type SettingUserResponse struct {
	Id    *uint                `json:"id"`
	Name  *string              `json:"name"`
	Email *string              `json:"email"`
	Type  *string              `json:"type"`
	Role  *SettingRoleResponse `json:"role"`
}

func SettingUserToResponse(user SettingUser, role *SettingRoleResponse) SettingUserResponse {
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
