package setting

import (
	"database/sql"
	"errors"
	"log"
	"server/database"
	"server/models"
	"server/models/setting"
	"server/utils"

	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

var (
	ErrRoleNotFound                    = errors.New("role not found")
	ErrCreateRoleWithFullPermission    = errors.New("unable to create role with full permission")
	ErrUnableToDeleteCurrentlyUsedRole = errors.New("unable to delete currently used role")
)

type SettingRoleService struct {
	DB *sql.DB
}

func (service *SettingRoleService) Roles(qp *utils.QueryParams) ([]setting.SettingRoleResponse, int, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_roles" AS (
		SELECT
			id, name, description
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
		log.Printf("%s\n", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s\n", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	roles := make([]setting.SettingRoleResponse, 0)
	lastRole := setting.SettingRole{}
	permissions := make([]setting.SettingPermission, 0)
	for rows.Next() {
		tmpRole := setting.SettingRole{}
		tmpPermission := setting.SettingPermission{}
		err := rows.Scan(&tmpRole.Id, &tmpRole.Name, &tmpRole.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			log.Printf("%s\n", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		if lastRole.Id != tmpRole.Id {
			if lastRole.Id != 0 {
				permissionsResponse := make([]setting.SettingPermissionResponse, 0)
				for _, permission := range permissions {
					permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
				}
				roles = append(roles, setting.SettingRoleToResponse(lastRole, permissionsResponse))
			}
			lastRole = tmpRole
			permissions = make([]setting.SettingPermission, 0)
		}
		permissions = append(permissions, tmpPermission)
	}
	if lastRole.Id != 0 {
		permissionsResponse := make([]setting.SettingPermissionResponse, 0)
		for _, permission := range permissions {
			permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
		}
		roles = append(roles, setting.SettingRoleToResponse(lastRole, permissionsResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "setting.role"`)
	qp.FilterIntoBqb(bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s\n", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		log.Printf("%s\n", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	return roles, total, 200, nil
}

func (service *SettingRoleService) Role(id string) (setting.SettingRoleResponse, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_roles" AS (
		SELECT
			id, name, description
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
		log.Printf("%s", err)
		return setting.SettingRoleResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s\n", err)
		return setting.SettingRoleResponse{}, 500, utils.ErrInternalServer
	}

	var role setting.SettingRole
	permissions := make([]setting.SettingPermission, 0)
	for rows.Next() {
		tmpPermission := setting.SettingPermission{}
		err := rows.Scan(&role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			log.Printf("%s\n", err)
			return setting.SettingRoleResponse{}, 500, utils.ErrInternalServer
		}

		permissions = append(permissions, tmpPermission)
	}

	if role.Id == 0 {
		return setting.SettingRoleResponse{}, 404, ErrRoleNotFound
	}

	permissionsResponse := make([]setting.SettingPermissionResponse, 0)
	for _, permission := range permissions {
		permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
	}

	return setting.SettingRoleToResponse(role, permissionsResponse), 200, nil
}

func (service *SettingRoleService) CreateRole(ctx *utils.CtxW, role *setting.SettingRoleCreateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForCreate(ctx.User.Id, ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`SELECT COUNT(*) FROM "setting.permission" WHERE id in (`)
	hasPermission := true
	for index, permissionId := range role.PermissionIds {
		if utils.ContainsString(utils.FULL_PERMISSION_IDS, utils.IntToStr(int(permissionId))) {
			for _, permission := range ctx.Permissions {
				if !utils.ContainsString([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCESS}, permission.Name) {
					hasPermission = false
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
		tx.Rollback()
		return 403, ErrCreateRoleWithFullPermission
	}

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	var count int
	err = tx.QueryRow(query, params...).Scan(&count)
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	if count != len(role.PermissionIds) {
		tx.Rollback()
		return 400, errors.New("permission_ids is invalid")
	}

	bqbQuery = bqb.New(`INSERT INTO "setting.role" (name, description, cid, ctime, mid, mtime) VALUES (?, ?, ?, ?, ?, ?) RETURNING id`, role.Name, role.Description, commonModel.CId, commonModel.CTime, commonModel.MId, commonModel.MTime)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	var roleId uint
	err = tx.QueryRow(query, params...).Scan(&roleId)
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
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
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	return 201, nil
}

func (service *SettingRoleService) UpdateRole(ctx *utils.CtxW, id string, role *setting.SettingRoleUpdateRequest) (int, error) {
	commonModel := models.CommonModel{}
	commonModel.PrepareForUpdate(ctx.User.Id)

	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`UPDATE "setting.role" SET mid = ?, mtime = ?`, commonModel.MId, commonModel.MTime)
	utils.PrepareUpdateBqbQuery(bqbQuery, role)
	bqbQuery.Space(`WHERE id = ?`, id)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		tx.Rollback()
		return 404, ErrRoleNotFound
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
			tx.Rollback()
			log.Printf("%s\n", err)
			return 500, utils.ErrInternalServer
		}

		_, err = tx.Exec(query, params...)
		if err != nil {
			tx.Rollback()
			log.Printf("%s\n", err)
			return 500, utils.ErrInternalServer
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
			tx.Rollback()
			log.Printf("%s\n", err)
			return 500, utils.ErrInternalServer
		}

		_, err = tx.Exec(query, params...)
		if err != nil {
			tx.Rollback()
			log.Printf("%s\n", err)
			return 500, utils.ErrInternalServer
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Printf("%s\n", err)
		return 500, utils.ErrInternalServer
	}

	return 200, nil
}

func (service *SettingRoleService) DeleteRole(id string) (int, error) {
	tx, err := service.DB.Begin()
	if err != nil {
		log.Printf("%v\n", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery := bqb.New(`DELETE FROM "setting.role_permission" WHERE setting_role_id = ?`, id)
	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v\n", err)
		return 500, utils.ErrInternalServer
	}

	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		log.Printf("%v\n", err)
		return 500, utils.ErrInternalServer
	}

	bqbQuery = bqb.New(`DELETE FROM "setting.role" WHERE id = ?`, id)
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("%v\n", err)
		return 500, utils.ErrInternalServer
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()

		switch err.(*pq.Error).Constraint {
		case database.FK_SETTING_ROLE_ID:
			return 409, ErrUnableToDeleteCurrentlyUsedRole
		}

		log.Printf("%v\n", err)
		return 500, utils.ErrInternalServer
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		tx.Rollback()
		return 404, ErrRoleNotFound
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("%v\n", err)
		return 500, utils.ErrInternalServer
	}

	return 200, nil
}
