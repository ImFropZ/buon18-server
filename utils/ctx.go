package utils

import (
	"system.buon18.com/m/models/setting"
)

type CtxKey struct{}

type CtxValue struct {
	User        *setting.SettingUser
	Role        *setting.SettingRole
	Permissions *[]setting.SettingPermission
}
