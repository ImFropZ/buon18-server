package setting

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
	"system.buon18.com/m/database"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/utils"
)

type SettingRoleService struct {
	DB *sql.DB
}

func (service *SettingRoleService) Roles(qp *utils.QueryParams) ([]setting.SettingRoleResponse, int, int, error) {
	bqbQuery := bqb.New(`WITH "limited_roles" AS (
	SELECT
		*
	FROM "setting.role"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.PaginationIntoBqb(bqbQuery)

	bqbQuery.Space(`)
	SELECT
		"limited_roles".id,
		"limited_roles".name,
		"limited_roles".description,
		"setting.permission".id,
		"setting.permission".name
	FROM "limited_roles"
	LEFT JOIN "setting.role_permission" ON "limited_roles".id = "setting.role_permission".setting_role_id
	LEFT JOIN "setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id`)

	qp.OrderByIntoBqb(bqbQuery, `"limited_roles".id ASC, "setting.permission".id ASC`)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	roles := make([]setting.SettingRoleResponse, 0)
	lastRole := setting.SettingRole{}
	permissions := make([]setting.SettingPermission, 0)
	for rows.Next() {
		tmpRole := setting.SettingRole{}
		tmpPermission := setting.SettingPermission{}
		err := rows.Scan(&tmpRole.Id, &tmpRole.Name, &tmpRole.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			slog.Error(fmt.Sprintf("%s\n", err))
			return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
		}

		if *lastRole.Id != *tmpRole.Id && lastRole.Id != nil {
			permissionsResponse := make([]setting.SettingPermissionResponse, 0)
			for _, permission := range permissions {
				permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
			}
			roles = append(roles, setting.SettingRoleToResponse(lastRole, &permissionsResponse))

			lastRole = tmpRole
			permissions = make([]setting.SettingPermission, 0)
		}

		permissions = append(permissions, tmpPermission)
	}
	if lastRole.Id != nil {
		permissionsResponse := make([]setting.SettingPermissionResponse, 0)
		for _, permission := range permissions {
			permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
		}
		roles = append(roles, setting.SettingRoleToResponse(lastRole, &permissionsResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "setting.role"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var total int
	if err := service.DB.QueryRow(query, params...).Scan(&total); err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return nil, 0, http.StatusInternalServerError, utils.ErrInternalServer
	}

	return roles, total, http.StatusOK, nil
}

func (service *SettingRoleService) Role(id string) (setting.SettingRoleResponse, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_roles" AS (
		SELECT
			*
		FROM "setting.role"
		WHERE id = ?
	)
	SELECT
		"limited_roles".id,
		"limited_roles".name,
		"limited_roles".description,
		"setting.permission".id,
		"setting.permission".name
	FROM "limited_roles"
	LEFT JOIN "setting.role_permission" ON "limited_roles".id = "setting.role_permission".setting_role_id
	LEFT JOIN "setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id
	ORDER BY "limited_roles".id ASC, "setting.permission".id ASC`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s", err))
		return setting.SettingRoleResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return setting.SettingRoleResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
	}

	var role setting.SettingRole
	permissions := make([]setting.SettingPermission, 0)
	for rows.Next() {
		tmpPermission := setting.SettingPermission{}
		err := rows.Scan(&role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			slog.Error(fmt.Sprintf("%s\n", err))
			return setting.SettingRoleResponse{}, http.StatusInternalServerError, utils.ErrInternalServer
		}

		permissions = append(permissions, tmpPermission)
	}

	if role.Id == nil {
		return setting.SettingRoleResponse{}, http.StatusNotFound, utils.ErrRoleNotFound
	}

	permissionsResponse := make([]setting.SettingPermissionResponse, 0)
	for _, permission := range permissions {
		permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
	}

	return setting.SettingRoleToResponse(role, &permissionsResponse), http.StatusOK, nil
}

func (service *SettingRoleService) CreateRole(ctx *utils.CtxValue, role *setting.SettingRoleCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(*ctx.User.Id, *ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`SELECT COUNT(*) FROM "setting.permission" WHERE id in (`)
	hasPermission := true
	for index, permissionId := range role.PermissionIds {
		if utils.ContainsString(utils.FULL_PERMISSION_IDS, utils.IntToStr(int(permissionId))) {
			for _, permission := range *ctx.Permissions {
				if permission.Name != nil {
					if !utils.ContainsString([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCESS}, *permission.Name) {
						hasPermission = false
					}
				}
			}
		}
		bqbQuery.Space(`?`, permissionId)
		if index != len(role.PermissionIds)-1 {
			bqbQuery.Space(`,`)
		}
	}
	bqbQuery.Space(`)`)

	if !hasPermission {
		return http.StatusForbidden, utils.ErrCreateRoleWithFullPermission
	}

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	var count int
	if err := tx.QueryRow(query, params...).Scan(&count); err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if count != len(role.PermissionIds) {
		return http.StatusBadRequest, utils.ErrPermissionNotFound
	}

	bqbQuery = bqb.New(`INSERT INTO "setting.role" (name, description, cid, ctime, mid, mtime) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`, role.Name, role.Description, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	var roleId uint
	if err := tx.QueryRow(query, params...).Scan(&roleId); err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`INSERT INTO "setting.role_permission" (setting_role_id, setting_permission_id, cid, ctime, mid, mtime) VALUES`)
	for index, permissionId := range role.PermissionIds {
		bqbQuery.Space(`(?, ?, ?, ?, ?, ?)`, roleId, permissionId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
		if index != len(role.PermissionIds)-1 {
			bqbQuery.Space(`,`)
		}
	}

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err := tx.Exec(query, params...); err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if err := tx.Commit(); err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusCreated, nil
}

func (service *SettingRoleService) UpdateRole(ctx *utils.CtxValue, id string, role *setting.SettingRoleUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForUpdate(*ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`UPDATE "setting.role" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	if role.Name != nil {
		bqbQuery.Space(`name = ?`, *role.Name)
	}
	if role.Description != nil {
		bqbQuery.Space(`description = ?`, *role.Description)
	}
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return http.StatusNotFound, utils.ErrRoleNotFound
	}

	if role.AddPermissionIds != nil {
		bqbQuery := bqb.New(`INSERT INTO "setting.role_permission" (setting_role_id, setting_permission_id, cid, ctime, mid, mtime) VALUES`)
		for index, permissionId := range *role.AddPermissionIds {
			bqbQuery.Space(`(?, ?, ?, ?, ?, ?)`, id, permissionId, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
			if index != len(*role.AddPermissionIds)-1 {
				bqbQuery.Space(`,`)
			}
		}

		query, params, err := bqbQuery.ToPgsql()
		if err != nil {
			slog.Error(fmt.Sprintf("%s\n", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}

		if _, err := tx.Exec(query, params...); err != nil {
			slog.Error(fmt.Sprintf("%s\n", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}
	}

	if role.RemovePermissionIds != nil {
		bqbQuery := bqb.New(`DELETE FROM "setting.role_permission" WHERE setting_role_id = ? AND setting_permission_id IN (`, id)
		for index, permissionId := range *role.RemovePermissionIds {
			bqbQuery.Space(`?`, permissionId)
			if index != len(*role.RemovePermissionIds)-1 {
				bqbQuery.Space(`,`)
			}
		}
		bqbQuery.Space(`)`)
		query, params, err := bqbQuery.ToPgsql()
		if err != nil {
			slog.Error(fmt.Sprintf("%s\n", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}

		if _, err := tx.Exec(query, params...); err != nil {
			slog.Error(fmt.Sprintf("%s\n", err))
			return http.StatusInternalServerError, utils.ErrInternalServer
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Error(fmt.Sprintf("%s\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}

func (service *SettingRoleService) DeleteRole(id string) (int, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		slog.Error(fmt.Sprintf("%v\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}
	defer tx.Rollback()

	bqbQuery := bqb.New(`DELETE FROM "setting.role_permission" WHERE setting_role_id = ?`, id)
	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if _, err := tx.Exec(query, params...); err != nil {
		slog.Error(fmt.Sprintf("%v\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`DELETE FROM "setting.role" WHERE id = ?`, id)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		slog.Error(fmt.Sprintf("%v\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		switch err.(*pq.Error).Constraint {
		case database.FK_SETTING_ROLE_ID:
			return http.StatusConflict, utils.ErrResourceInUsed
		}

		slog.Error(fmt.Sprintf("%v\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return http.StatusNotFound, utils.ErrRoleNotFound
	}

	if err := tx.Commit(); err != nil {
		slog.Error(fmt.Sprintf("%v\n", err))
		return http.StatusInternalServerError, utils.ErrInternalServer
	}

	return http.StatusOK, nil
}
