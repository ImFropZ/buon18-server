package utils

import (
	"errors"
	"log"
	"server/models"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoUserCtx        = errors.New("unable to get user context")
	ErrNoRoleCtx        = errors.New("unable to get role context")
	ErrNoPermissionsCtx = errors.New("unable to get permissions context")
)

type CtxW struct {
	User        models.SettingUser
	Role        models.SettingRole
	Permissions []models.SettingPermission
}

func Ctx(c *gin.Context) (CtxW, error) {
	// -- Get user
	var ctxUser models.SettingUser
	if cUser, err := c.Get("user"); !err {
		log.Printf("Error getting user id: %v\n", err)
		return CtxW{}, ErrNoUserCtx
	} else {
		ctxUser = cUser.(models.SettingUser)
	}

	// -- Get role
	var ctxRole models.SettingRole
	if cRole, err := c.Get("role"); !err {
		log.Printf("Error getting user id: %v\n", err)
		return CtxW{}, ErrNoRoleCtx
	} else {
		ctxRole = cRole.(models.SettingRole)
	}

	// -- Get permissions
	var ctxPermissions []models.SettingPermission
	if cPermissions, err := c.Get("permissions"); !err {
		log.Printf("Error getting user id: %v\n", err)
		return CtxW{}, ErrNoPermissionsCtx
	} else {
		ctxPermissions = cPermissions.([]models.SettingPermission)
	}

	return CtxW{
		User:        ctxUser,
		Role:        ctxRole,
		Permissions: ctxPermissions,
	}, nil
}
