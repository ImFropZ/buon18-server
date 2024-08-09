package utils

import (
	"errors"
	"log"
	"server/models/setting"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoUserCtx        = errors.New("unable to get user context")
	ErrNoRoleCtx        = errors.New("unable to get role context")
	ErrNoPermissionsCtx = errors.New("unable to get permissions context")
)

type CtxW struct {
	User        setting.SettingUser
	Role        setting.SettingRole
	Permissions []setting.SettingPermission
}

func Ctx(c *gin.Context) (CtxW, error) {
	// -- Get user
	var ctxUser setting.SettingUser
	if cUser, err := c.Get("user"); !err {
		log.Printf("Error getting user id: %v\n", err)
		return CtxW{}, ErrNoUserCtx
	} else {
		ctxUser = cUser.(setting.SettingUser)
	}

	// -- Get role
	var ctxRole setting.SettingRole
	if cRole, err := c.Get("role"); !err {
		log.Printf("Error getting user id: %v\n", err)
		return CtxW{}, ErrNoRoleCtx
	} else {
		ctxRole = cRole.(setting.SettingRole)
	}

	// -- Get permissions
	var ctxPermissions []setting.SettingPermission
	if cPermissions, err := c.Get("permissions"); !err {
		log.Printf("Error getting user id: %v\n", err)
		return CtxW{}, ErrNoPermissionsCtx
	} else {
		ctxPermissions = cPermissions.([]setting.SettingPermission)
	}

	return CtxW{
		User:        ctxUser,
		Role:        ctxRole,
		Permissions: ctxPermissions,
	}, nil
}
