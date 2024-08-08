package middlewares

import (
	"database/sql"
	"errors"
	"log"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nullism/bqb"
)

func Authenticate(DB *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// -- Get token
		token, err := utils.RemoveBearer(c.GetHeader("Authorization"))
		if err != nil {
			c.JSON(401, utils.NewErrorResponse(401, "token is required"))
			c.Abort()
			return
		}

		// -- Validate token
		claims, err := utils.ValidateWebToken(token)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(401, utils.NewErrorResponse(401, "token expired"))
				c.Abort()
				return
			}

			c.JSON(401, utils.NewErrorResponse(401, "invalid token"))
			c.Abort()
			return
		}

		// -- Prepare sql query
		query, params, err := bqb.New(`
		SELECT 
			"setting.user".id, 
			"setting.user".name, 
			"setting.user".email, 
			"setting.user".typ, 
			COALESCE("setting.role".id, 0), 
			COALESCE("setting.role".name, ''), 
			COALESCE("setting.role".description, ''), 
			COALESCE("setting.permission".id, 0), 
			COALESCE("setting.permission".name, '')
		FROM 
			"setting.user"
		LEFT JOIN 
			"setting.role" ON "setting.user".setting_role_id = "setting.role".id
		LEFT JOIN 
			"setting.role_permission" ON "setting.role".id = "setting.role_permission".setting_role_id 
		LEFT JOIN 
			"setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id
		WHERE 
			"setting.user".email = ?
		ORDER BY "setting.user".email, "setting.role".id, "setting.permission".id`, claims.Email).ToPgsql()
		if err != nil {
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			c.Abort()
			return
		}

		// -- Validate user
		var user models.SettingUser
		var role models.SettingRole
		permissions := make([]models.SettingPermission, 0)
		if row := DB.QueryRow(query, params...); row.Err() != nil {
			log.Printf("Error querying user: %v\n", row.Err())
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		} else {
			var tmpPermission models.SettingPermission
			if err := row.Scan(&user.Id, &user.Name, &user.Email, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name); err != nil {
				if err == sql.ErrNoRows {
					c.JSON(401, utils.NewErrorResponse(401, "contact your administrator to create an account"))
					c.Abort()
					return
				}

				log.Printf("Error scanning user: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				c.Abort()
				return
			}

			// -- Append permission
			permissions = append(permissions, tmpPermission)
		}

		// -- Set user info
		c.Set("user", user)
		c.Set("role", role)
		c.Set("permissions", permissions)

		log.Printf("Permissions: %v\n", permissions)

		c.Next()
	}
}
