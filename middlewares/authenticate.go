package middlewares

import (
	"database/sql"
	"log"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
)

func Authenticate(DB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// -- Get token from header
		token := c.GetHeader("Authorization")

		// -- Remove Bearer schema
		token, err := utils.RemoveBearer(token)
		if err != nil {
			c.JSON(401, utils.NewErrorResponse(401, "token is required"))
			c.Abort()
			return
		}

		// -- Validate token
		claims, err := utils.ValidateWebToken(token)
		if err != nil {
			c.JSON(401, utils.NewErrorResponse(401, "invalid token"))
			c.Abort()
			return
		}

		// -- Prepare sql query
		query, params, err := bqb.New("SELECT id, role FROM \"user\" WHERE email = ?", claims.Email).ToPgsql()
		if err != nil {
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			c.Abort()
			return
		}

		// -- Query user from database
		var User models.User
		if row := DB.QueryRow(query, params...); row.Err() != nil {
			log.Printf("Error querying user : %v\n", row.Err())
			c.JSON(401, utils.NewErrorResponse(401, "user not found"))
			c.Abort()
			return
		} else {
			if err := row.Scan(&User.Id, &User.Role); err != nil {
				log.Printf("Error scanning user: %v\n", err)
				c.JSON(401, utils.NewErrorResponse(401, "user not found"))
				c.Abort()
				return
			}
		}

		// -- Set user info
		c.Set("role", User.Role)
		c.Set("user_id", User.Id)

		c.Next()
	}
}
