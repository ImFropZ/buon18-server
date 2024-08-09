package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"server/models"
	"server/models/setting"
	"server/utils"

	"github.com/nullism/bqb"
)

var (
	ErrAccountNotFound        = errors.New("contact your administrator to create an account")
	ErrInvalidEmailOrPassword = errors.New("invalid email or password")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrInvalidOldPassword     = errors.New("invalid old password")
	ErrUserNotFound           = errors.New("user not found")
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

type AuthService struct {
	DB *sql.DB
}

func (service *AuthService) Me(ctx *utils.CtxW) (setting.SettingUserResponse, int, error) {
	return setting.SettingUserToResponse(ctx.User, ctx.Role, ctx.Permissions), 200, nil
}

func (service *AuthService) Login(loginRequest *LoginRequest) (models.TokenAndRefreshToken, int, error) {
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
	ORDER BY "setting.user".email, "setting.role".id, "setting.permission".id`, loginRequest.Email).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	}

	// -- Validate user
	var user setting.SettingUser
	var role setting.SettingRole
	permissions := make([]setting.SettingPermission, 0)
	if row := service.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	} else {
		var tmpPermission setting.SettingPermission
		if err := row.Scan(&user.Email, &user.Pwd, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name); err != nil {
			if err == sql.ErrNoRows {
				return models.TokenAndRefreshToken{}, 401, ErrAccountNotFound
			}

			log.Printf("Error scanning user: %v\n", err)
			return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
		}

		// -- Append permission
		permissions = append(permissions, tmpPermission)
	}

	if user.Email != loginRequest.Email || (!utils.ComparePwd(loginRequest.Password, user.Pwd) && user.Pwd != "") {
		return models.TokenAndRefreshToken{}, 401, ErrInvalidEmailOrPassword
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
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	}

	return models.TokenAndRefreshToken{
		Token:        token,
		RefreshToken: refreshToken,
	}, 200, nil
}

func (service *AuthService) RefreshToken(refreshTokenRequest *RefreshTokenRequest) (models.TokenAndRefreshToken, int, error) {
	refreshClaims, refreshErr := utils.ValidateRefreshToken(refreshTokenRequest.RefreshToken)
	if refreshErr != nil {
		return models.TokenAndRefreshToken{}, 401, ErrInvalidRefreshToken
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
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	}

	// -- Validate user
	var user setting.SettingUser
	var role setting.SettingRole
	permissions := make([]setting.SettingPermission, 0)
	if row := service.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	} else {
		var tmpPermission setting.SettingPermission
		if err := row.Scan(&user.Email, &user.Pwd, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name); err != nil {
			if err == sql.ErrNoRows {
				return models.TokenAndRefreshToken{}, 401, ErrAccountNotFound
			}

			log.Printf("Error scanning user: %v\n", err)
			return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
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
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		return models.TokenAndRefreshToken{}, 500, utils.ErrInternalServer
	}

	return models.TokenAndRefreshToken{
		Token:        token,
		RefreshToken: refreshToken,
	}, 200, nil
}

func (service *AuthService) UpdatePassword(ctx *utils.CtxW, updatePasswordRequest *UpdatePasswordRequest) (string, int, error) {
	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT COALESCE(pwd, '') FROM "setting.user" WHERE id = ?`, ctx.User.Id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v\n", err)
		return "", 500, utils.ErrInternalServer
	}

	// -- Get user from db
	var user setting.SettingUser
	if row := service.DB.QueryRow(query, params...); row.Err() != nil {
		log.Printf("Error querying user: %v\n", row.Err())
		return "", 500, utils.ErrInternalServer
	} else {
		if err := row.Scan(&user.Pwd); err != nil {
			log.Printf("Error scanning user: %v\n", err)
			return "", 500, utils.ErrInternalServer
		}
	}

	// -- Validate old password if exists
	if user.Pwd != "" {
		if ok := utils.ComparePwd(updatePasswordRequest.OldPassword, user.Pwd); !ok {
			return "", 401, ErrInvalidOldPassword
		}
	}

	// -- Update pwd
	if hashedPwd, err := utils.HashPwd(updatePasswordRequest.NewPassword); err == nil {
		// -- Begin transaction
		tx, err := service.DB.Begin()
		if err != nil {
			log.Printf("Error beginning transaction: %v\n", err)
			return "", 500, utils.ErrInternalServer
		}

		// -- Prepare sql query
		query, params, err = bqb.New(`UPDATE "setting.user" SET pwd = ? WHERE id = ?`, hashedPwd, ctx.User.Id).ToPgsql()
		if err != nil {
			log.Printf("Error preparing query: %v\n", err)
			return "", 500, utils.ErrInternalServer
		}

		// -- Update pwd
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()

			log.Printf("Error updating password: %v\n", err)
			return "", 500, utils.ErrInternalServer
		}

		// -- Commit transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v\n", err)
			return "", 500, utils.ErrInternalServer
		}

		return "success", 200, nil
	} else {
		log.Printf("Error hashing password: %v\n", err)
		return "", 500, utils.ErrInternalServer
	}
}

func (service *AuthService) UpdateProfile(ctx *utils.CtxW, updateData *UpdateProfileRequest) (string, int, error) {
	if ctx.User.Typ != models.SettingUserTypUser {
		return "", 403, utils.ErrForbidden
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
		return "", 500, utils.ErrInternalServer
	}

	// -- Begin transaction
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		return "", 500, utils.ErrInternalServer
	}

	// -- Check if user exists
	var count int
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error querying user: %v\n", row.Err())
		return "", 500, utils.ErrInternalServer
	} else {
		if err := row.Scan(&count); err != nil {
			tx.Rollback()
			log.Printf("Error scanning user: %v\n", err)
			return "", 500, utils.ErrInternalServer
		}
	}
	if count == 0 {
		tx.Rollback()
		return "", 404, ErrUserNotFound
	}

	// -- Loop through request fields
	updateFeilds := map[string]string{}
	v := reflect.ValueOf(*updateData)
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
					return "", 403, utils.ErrForbidden
				}
				updateFeilds[fieldName] = *field.Interface().(*string)
			default:
				return "", 400, utils.ErrBadRequest
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
		return "", 500, utils.ErrInternalServer
	}

	// -- Update user
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error updating user: %v\n", err)
		return "", 500, utils.ErrInternalServer
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		return "", 500, utils.ErrInternalServer
	}

	return "success", 200, nil
}
