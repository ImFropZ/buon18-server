package setting

import (
	"database/sql"
	"log"
	"server/models"
	"server/utils"

	"github.com/nullism/bqb"
)

type SettingRoleService struct {
	DB *sql.DB
}

func (service *SettingRoleService) Roles(qp *utils.QueryParams) ([]models.SettingRoleResponse, int, int, error) {
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

	roles := make([]models.SettingRoleResponse, 0)
	lastRole := models.SettingRole{}
	permissions := make([]models.SettingPermission, 0)
	for rows.Next() {
		tmpRole := models.SettingRole{}
		tmpPermission := models.SettingPermission{}
		err := rows.Scan(&tmpRole.Id, &tmpRole.Name, &tmpRole.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			log.Printf("%s\n", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		if lastRole.Id != tmpRole.Id {
			if lastRole.Id != 0 {
				roles = append(roles, models.SettingRoleToResponse(lastRole, permissions))
			}
			lastRole = tmpRole
			permissions = make([]models.SettingPermission, 0)
		}
		permissions = append(permissions, tmpPermission)
	}
	if lastRole.Id != 0 {
		roles = append(roles, models.SettingRoleToResponse(lastRole, permissions))
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

	return roles, total, 0, nil
}
