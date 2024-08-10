package setting

import (
	"database/sql"
	"errors"
	"log"
	"server/models/setting"
	"server/utils"

	"github.com/nullism/bqb"
)

var (
	ErrRoleNotFound = errors.New("role not found")
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
