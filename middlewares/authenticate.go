package middlewares

import (
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// -- Get token from header
		token := c.GetHeader("Authorization")

		// -- Remove Bearer schema
		token, err := utils.RemoveBearer(token)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// -- Validate token
		claims, err := utils.ValidateWebToken(token)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// -- Set user info
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}
