package middlewares

import (
	"log"
	"server/models/setting"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorize(allowPermissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// -- Get role
		var ctxPermissions []setting.SettingPermission
		if cPermissions, ok := c.Get("permissions"); ok {
			ctxPermissions = cPermissions.([]setting.SettingPermission)
		} else {
			log.Printf("permission not found in context\n")
			c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
			c.Abort()
			return
		}

		// -- Add FULL_ACCESS permission
		allowPermissions = append(allowPermissions, "FULL_ACCESS")

		allow := false
		for _, permission := range allowPermissions {
			// -- Check permission
			for _, ctxPermission := range ctxPermissions {
				if strings.EqualFold(ctxPermission.Name, permission) {
					allow = true
					break
				}
			}
		}

		if !allow {
			c.JSON(403, utils.NewErrorResponse(403, utils.ErrForbidden.Error()))
			c.Abort()
			return
		}

		c.Next()
	}
}
