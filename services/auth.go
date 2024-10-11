package services

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/nullism/bqb"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/utils"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword *string `json:"old_password" validate:"required"`
	NewPassword *string `json:"new_password" validate:"required"`
}

type UpdateProfileRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	RoleId *uint   `json:"role_id"`
}

type AuthService struct {
	DB *sql.DB
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
		slog.Error(fmt.Sprintf("Error preparing query: %v\n", err))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Validate user
	var user setting.SettingUser
	var role setting.SettingRole
	permissions := make([]setting.SettingPermission, 0)
	if row := service.DB.QueryRow(query, params...); row.Err() != nil {
		slog.Error(fmt.Sprintf("Error querying user: %v\n", row.Err()))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	} else {
		var tmpPermission setting.SettingPermission
		if err := row.Scan(&user.Email, &user.Pwd, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name); err != nil {
			if err == sql.ErrNoRows {
				return models.TokenAndRefreshToken{}, http.StatusNotFound, utils.ErrUserAccountNotFound
			}

			slog.Error(fmt.Sprintf("Error scanning user: %v\n", err))
			return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
		}

		// -- Append permission
		permissions = append(permissions, tmpPermission)
	}

	if (user.Email != nil && *user.Email != loginRequest.Email) || (user.Pwd != nil && !utils.ComparePwd(loginRequest.Password, *user.Pwd) && *user.Pwd != "") {
		return models.TokenAndRefreshToken{}, http.StatusBadRequest, utils.ErrInvalidEmailOrPassword
	}

	// -- Generate token
	permissionNames := make([]string, 0)
	for _, value := range permissions {
		if value.Name != nil {
			permissionNames = append(permissionNames, *value.Name)
		}
	}

	if user.Email == nil || role.Name == nil {
		return models.TokenAndRefreshToken{}, http.StatusNotFound, utils.ErrUserAccountNotFound
	}
	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       *user.Email,
		Role:        *role.Name,
		Permissions: permissionNames,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating web token: %v\n", err))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: *user.Email,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating refresh token: %v\n", err))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return models.TokenAndRefreshToken{
		Token:        token,
		RefreshToken: refreshToken,
	}, http.StatusOK, nil
}

func (service *AuthService) RefreshToken(refreshTokenRequest *RefreshTokenRequest) (models.TokenAndRefreshToken, int, error) {
	refreshClaims, refreshErr := utils.ValidateRefreshToken(refreshTokenRequest.RefreshToken)
	if refreshErr != nil {
		return models.TokenAndRefreshToken{}, http.StatusBadRequest, utils.ErrInvalidRefreshToken
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`
	SELECT
		"setting.user".email,
		"setting.user".pwd,
		"setting.user".typ,
		"setting.role".id
		"setting.role".name
		"setting.role".description
		"setting.permission".id
		"setting.permission".name
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
		slog.Error(fmt.Sprintf("Error preparing query: %v\n", err))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Validate user
	var user setting.SettingUser
	var role setting.SettingRole
	permissions := make([]setting.SettingPermission, 0)
	if row := service.DB.QueryRow(query, params...); row.Err() != nil {
		slog.Error(fmt.Sprintf("Error querying user: %v\n", row.Err()))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	} else {
		var tmpPermission setting.SettingPermission
		if err := row.Scan(&user.Email, &user.Pwd, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name); err != nil {
			if err == sql.ErrNoRows {
				return models.TokenAndRefreshToken{}, http.StatusNotFound, utils.ErrUserAccountNotFound
			}

			slog.Error(fmt.Sprintf("Error scanning user: %v\n", err))
			return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
		}

		// -- Append permission
		permissions = append(permissions, tmpPermission)
	}

	// -- Generate token
	permissionNames := make([]string, 0)
	for _, value := range permissions {
		if value.Name != nil {
			permissionNames = append(permissionNames, *value.Name)
		}
	}

	if user.Email == nil || role.Name == nil {
		return models.TokenAndRefreshToken{}, http.StatusNotFound, utils.ErrUserAccountNotFound
	}
	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email:       *user.Email,
		Role:        *role.Name,
		Permissions: permissionNames,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating web token: %v\n", err))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: *user.Email,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("Error generating refresh token: %v\n", err))
		return models.TokenAndRefreshToken{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return models.TokenAndRefreshToken{
		Token:        token,
		RefreshToken: refreshToken,
	}, http.StatusOK, nil
}

func (service *AuthService) UpdatePassword(ctx *utils.CtxValue, updatePasswordRequest *UpdatePasswordRequest) (string, int, error) {
	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT COALESCE(pwd, '') FROM "setting.user" WHERE id = ?`, ctx.User.Id).ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("Error preparing query: %v\n", err))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Get user from db
	var user setting.SettingUser
	if row := service.DB.QueryRow(query, params...); row.Err() != nil {
		slog.Error(fmt.Sprintf("Error querying user: %v\n", row.Err()))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	} else {
		if err := row.Scan(&user.Pwd); err != nil {
			slog.Error(fmt.Sprintf("Error scanning user: %v\n", err))
			return "", http.StatusInternalServerError, utils.ErrInternalServer
		}
	}

	// -- Validate old password if exists
	if (user.Pwd != nil && *user.Pwd != "") && updatePasswordRequest.OldPassword != nil {
		if ok := utils.ComparePwd(*updatePasswordRequest.OldPassword, *user.Pwd); !ok {
			return "", http.StatusBadRequest, utils.ErrInvalidOldPassword
		}
	}

	// -- Update pwd
	if hashedPwd, err := utils.HashPwd(*updatePasswordRequest.NewPassword); err == nil {
		// -- Begin transaction
		tx, err := service.DB.Begin()
		if err != nil {
			slog.Error(fmt.Sprintf("Error beginning transaction: %v\n", err))
			return "", http.StatusInternalServerError, utils.ErrInternalServer
		}
		defer tx.Rollback()

		// -- Prepare sql query
		query, params, err = bqb.New(`UPDATE "setting.user" SET pwd = ? WHERE id = ?`, hashedPwd, ctx.User.Id).ToPgsql()
		if err != nil {
			slog.Error(fmt.Sprintf("Error preparing query: %v\n", err))
			return "", http.StatusInternalServerError, utils.ErrInternalServer
		}

		// -- Update pwd
		if _, err := tx.Exec(query, params...); err != nil {
			slog.Error(fmt.Sprintf("Error updating password: %v\n", err))
			return "", http.StatusInternalServerError, utils.ErrInternalServer
		}

		// -- Commit transaction
		if err := tx.Commit(); err != nil {
			slog.Error(fmt.Sprintf("Error committing transaction: %v\n", err))
			return "", http.StatusInternalServerError, utils.ErrInternalServer
		}

		return "success", http.StatusOK, nil
	} else {
		slog.Error(fmt.Sprintf("Error hashing password: %v\n", err))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	}
}

func (service *AuthService) UpdateProfile(ctx *utils.CtxValue, updateData *UpdateProfileRequest) (string, int, error) {
	// -- Prepare sql query
	query, params, err := bqb.New(`
	SELECT
		count(*)
	FROM
		"setting.user"
	WHERE
		id = ?`, ctx.User.Id).ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("Error preparing query: %v\n", err))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Begin transaction
	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("Error beginning transaction: %v\n", err))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	// -- Check if user exists
	var count int
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		slog.Error(fmt.Sprintf("Error querying user: %v\n", row.Err()))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	} else {
		if err := row.Scan(&count); err != nil {
			slog.Error(fmt.Sprintf("Error scanning user: %v\n", err))
			return "", http.StatusInternalServerError, utils.ErrInternalServer
		}
	}
	if count == 0 {
		return "", http.StatusNotFound, utils.ErrUserNotFound
	}

	// -- Prepare sql query
	bqbQuery := bqb.New(`UPDATE "setting.user" SET`)
	if updateData.Name != nil {
		bqbQuery = bqbQuery.Space(`name = ?`, *updateData.Name)
	}
	if updateData.Email != nil {
		bqbQuery = bqbQuery.Space(`email = ?`, *updateData.Email)
	}
	if updateData.RoleId != nil {
		// -- Check if user is allowed to change role
		permissionArr := []string{}
		for _, permission := range *ctx.Permissions {
			if permission.Name != nil {
				permissionArr = append(permissionArr, *permission.Name)
			}
		}
		if utils.ContainsString(permissionArr, "FULL_ACCESS") {
			tx.Rollback()
			return "", http.StatusForbidden, utils.ErrForbidden
		}

		bqbQuery = bqbQuery.Space(`setting_role_id = ?`, *updateData.RoleId)
	}

	query, params, err = bqbQuery.Space(`WHERE id = ?`, ctx.User.Id).ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("Error preparing query: %v\n", err))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Update user
	if _, err := tx.Exec(query, params...); err != nil {
		slog.Error(fmt.Sprintf("Error updating user: %v\n", err))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		slog.Error(fmt.Sprintf("Error committing transaction: %v\n", err))
		return "", http.StatusInternalServerError, utils.ErrInternalServer
	}

	return "success", http.StatusOK, nil
}
