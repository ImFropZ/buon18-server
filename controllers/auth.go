package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UpdateProfileRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	RoleId *uint   `json:"role_id"`
}

type AuthHandler struct {
	DB *sql.DB
}

func (handler *AuthHandler) Me(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"data": models.SettingUserToResponse(ctx.User, ctx.Role, ctx.Permissions),
	}))
}

func (handler *AuthHandler) Login(c *gin.Context) {
	// -- Parse request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`
	SELECT 
		"setting.user".email, 
		COALESCE("setting.user".pwd, ''), 
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
	ORDER BY "setting.user".email, "setting.role".id, "setting.permission".id`, req.Email).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Validate user
	var user models.SettingUser
	var role models.SettingRole
	permissions := make([]models.SettingPermission, 0)
	if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		var tmpPermission models.SettingPermission
		if err := row.Scan(&user.Email, &user.Pwd, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(401, utils.NewErrorResponse(401, "contact your administrator to create an account"))
				return
			}

			log.Printf("Error scanning user: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Append permission
		permissions = append(permissions, tmpPermission)
	}

	if user.Email != req.Email || (!utils.ComparePwd(req.Password, user.Pwd) && user.Pwd != "") {
		c.JSON(401, utils.NewErrorResponse(401, "Invalid email or password"))
		return
	}

	// -- Generate token
	permissionNames := make([]string, 0)
	for _, value := range permissions {
		permissionNames = append(permissionNames, value.Name)
	}
	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       user.Email,
		Role:        role.Name,
		Permissions: permissionNames,
	})
	if err != nil {
		log.Printf("Error generating web token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	}))
}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	// -- Parse request
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain refresh_token field"))
		return
	}

	refreshClaims, refreshErr := utils.ValidateRefreshToken(req.RefreshToken)
	if refreshErr != nil {
		c.JSON(401, utils.NewErrorResponse(401, "invalid refresh token"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`
	SELECT 
		"setting.user".email, 
		COALESCE("setting.user".pwd, ''), 
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
	ORDER BY "setting.user".email, "setting.role".id, "setting.permission".id`, refreshClaims.Email).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Validate user
	var user models.SettingUser
	var role models.SettingRole
	permissions := make([]models.SettingPermission, 0)
	if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		var tmpPermission models.SettingPermission
		if err := row.Scan(&user.Email, &user.Pwd, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(401, utils.NewErrorResponse(401, "contact your administrator to create an account"))
				return
			}

			log.Printf("Error scanning user: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Append permission
		permissions = append(permissions, tmpPermission)
	}

	// -- Generate token
	permissionNames := make([]string, 0)
	for _, value := range permissions {
		permissionNames = append(permissionNames, value.Name)
	}
	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       user.Email,
		Role:        role.Name,
		Permissions: permissionNames,
	})
	if err != nil {
		log.Printf("Error generating web token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	}))
}

func (handler *AuthHandler) UpdatePassword(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Parse request
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain old_password and new_password fields"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT COALESCE(pwd, '') FROM "setting.user" WHERE id = ?`, ctx.User.Id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get user from db
	var user models.SettingUser
	if row := handler.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&user.Pwd); err != nil {
			log.Printf("Error scanning user: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Validate old password if exists
	if user.Pwd != "" {
		if ok := utils.ComparePwd(req.OldPassword, user.Pwd); !ok {
			c.JSON(400, utils.NewErrorResponse(400, "invalid old password"))
			return
		}
	}

	// -- Update pwd
	if hashedPwd, err := utils.HashPwd(req.NewPassword); err == nil {
		// -- Begin transaction
		tx, err := handler.DB.Begin()
		if err != nil {
			log.Printf("Error beginning transaction: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Prepare sql query
		query, params, err = bqb.New(`UPDATE "setting.user" SET pwd = ? WHERE id = ?`, hashedPwd, ctx.User.Id).ToPgsql()
		if err != nil {
			log.Printf("Error preparing query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Update pwd
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()

			log.Printf("Error updating password: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Commit transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		c.JSON(200, utils.NewResponse(200, "success", nil))
	} else {
		log.Printf("Error hashing password: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
	}
}

func (handler *AuthHandler) UpdateProfile(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	if ctx.User.Typ != models.SettingUserTypUser {
		c.JSON(403, utils.NewErrorResponse(403, "forbidden"))
		return
	}

	// -- Parse request
	var updateData UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request"))
		return
	}
	if utils.IsAllFieldsNil(&updateData) {
		c.JSON(400, utils.NewErrorResponse(400, "request body should contain at least one field"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`
	SELECT 
		count(*) 
	FROM 
		"setting.user" 
	WHERE
		id = ?`, ctx.User.Id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Begin transaction
	tx, err := handler.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if user exists
	var count int
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&count); err != nil {
			tx.Rollback()
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}
	if count == 0 {
		tx.Rollback()
		c.JSON(404, utils.NewErrorResponse(404, "user not found"))
		return
	}

	// -- Loop through request fields
	updateFeilds := map[string]string{}
	v := reflect.ValueOf(updateData)
	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.IsNil() {
				continue
			}

			fieldName := utils.PascalToSnake(v.Type().Field(i).Name)
			switch fieldName {
			case "name":
				updateFeilds[fieldName] = *field.Interface().(*string)
			case "email":
				updateFeilds[fieldName] = *field.Interface().(*string)
			case "role_id":
				// -- Check if user is allowed to change role
				permissionArr := []string{}
				for _, permission := range ctx.Permissions {
					permissionArr = append(permissionArr, permission.Name)
				}
				if utils.ContainsString(permissionArr, "FULL_ACCESS") {
					tx.Rollback()
					c.JSON(403, utils.NewErrorResponse(403, "forbidden"))
					return
				}
				updateFeilds[fieldName] = *field.Interface().(*string)
			default:
				c.JSON(400, utils.NewErrorResponse(400, "invalid field"))
				return
			}
		}
	}

	// -- Prepare sql query
	bqbQuery := bqb.New(`UPDATE "setting.user" SET`)
	for key, value := range updateFeilds {
		bqbQuery = bqbQuery.Space(fmt.Sprintf(`%s = ?`, key), value)
	}
	query, params, err = bqbQuery.Space(`WHERE id = ?`, ctx.User.Id).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update user
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error updating user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", nil))
}
