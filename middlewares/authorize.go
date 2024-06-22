package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Role int

const (
	Admin Role = iota
	Editor
	User
)

func Authorize(role Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		// -- Get role
		userRole := c.GetString("role")

		// -- Compare role
		if convertRole(userRole) > role {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func convertRole(roleStr string) Role {
	switch strings.ToLower(roleStr) {
	case "admin":
		return Admin
	case "editor":
		return Editor
	case "user":
		return User
	default:
		return -1
	}
}
