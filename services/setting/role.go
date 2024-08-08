package setting

import (
	"database/sql"
	"fmt"
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

	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			bqbQuery.Space(fmt.Sprintf(`%s %s ?`, filter.Field, utils.MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]), filter.Value)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space("AND")
			}
		}
	}

	bqbQuery.Space(`OFFSET ? LIMIT ? )
	SELECT 
		"limited_roles".id,
		"limited_roles".name,
		"limited_roles".description,
		"setting.permission".id,
		"setting.permission".name
	FROM "limited_roles"
	LEFT JOIN "setting.role_permission" ON "limited_roles".id = "setting.role_permission".setting_role_id
	LEFT JOIN "setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id`, qp.Pagination.Offset, qp.Pagination.Limit)

	if len(qp.OrderBy) > 0 {
		bqbQuery.Space("ORDER BY")
		for index, sort := range qp.OrderBy {
			bqbQuery.Space(sort)
			if index < len(qp.OrderBy)-1 {
				bqbQuery.Space(",")
			}
		}
	} else {
		bqbQuery.Space(`ORDER BY "limited_roles".id ASC, "setting.permission".id ASC`)
	}

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
	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			bqbQuery.Space(fmt.Sprintf(`%s %s ?`, filter.Field, utils.MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]), filter.Value)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space("AND")
			}
		}
	}

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
