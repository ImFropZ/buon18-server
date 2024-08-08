package setting

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"server/models"
	"server/utils"

	"github.com/nullism/bqb"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type SettingUserService struct {
	DB *sql.DB
}

func (service *SettingUserService) Users(qp *utils.QueryParams) ([]models.SettingUserResponse, int, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_users" AS (
		SELECT 
			id, name, email, typ, setting_role_id
		FROM "setting.user"`)

	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			bqbQuery.Space(fmt.Sprintf("%s %s ?", filter.Field, utils.MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]), filter.Value)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space("AND")
			}
		}
	}

	bqbQuery.Space(`
		OFFSET ? LIMIT ?
	)
	SELECT 
		"limited_users".id, 
		"limited_users".name, 
		"limited_users".email, 
		"limited_users".typ,
		COALESCE("setting.role".id, 0), 
		COALESCE("setting.role".name, ''), 
		COALESCE("setting.role".description, ''), 
		COALESCE("setting.permission".id, 0), 
		COALESCE("setting.permission".name, '')
	FROM "limited_users"
	LEFT JOIN "setting.role" ON "limited_users".setting_role_id = "setting.role".id
	LEFT JOIN "setting.role_permission" ON "setting.role".id = "setting.role_permission".setting_role_id
	LEFT JOIN "setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id`, qp.Pagination.Offset, qp.Pagination.Limit)

	if len(qp.OrderBy) > 0 {
		bqbQuery.Space("ORDER BY")
		for index, sort := range qp.OrderBy {
			bqbQuery.Space(sort)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space(",")
			}
		}
		bqbQuery.Space(`, "setting.role".id ASC, "setting.permission".id ASC`)
	} else {
		bqbQuery.Space(`ORDER BY "limited_users".id ASC, "setting.role".id ASC, "setting.permission".id ASC`)
	}

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	usersResponse := make([]models.SettingUserResponse, 0)
	permission := make([]models.SettingPermission, 0)
	var lastUser models.SettingUser
	var lastRole models.SettingRole
	for rows.Next() {
		var tmpUser models.SettingUser
		var tmpRole models.SettingRole
		var tmpPermission models.SettingPermission
		err := rows.Scan(&tmpUser.Id, &tmpUser.Name, &tmpUser.Email, &tmpUser.Typ, &tmpRole.Id, &tmpRole.Name, &tmpRole.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			log.Printf("%s", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		if lastUser.Id != tmpUser.Id && lastUser.Id != 0 {
			usersResponse = append(usersResponse, models.SettingUserToResponse(lastUser, lastRole, permission))
			lastUser = tmpUser
			lastRole = tmpRole
			permission = make([]models.SettingPermission, 0)
			permission = append(permission, tmpPermission)
			continue
		}

		if lastUser.Id == 0 {
			lastUser = tmpUser
			lastRole = tmpRole
		}

		if tmpPermission.Id != 0 {
			permission = append(permission, tmpPermission)
		}
	}
	if lastUser.Id != 0 {
		usersResponse = append(usersResponse, models.SettingUserToResponse(lastUser, lastRole, permission))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "setting.user"`)
	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			bqbQuery.Space(fmt.Sprintf("%s %s ?", filter.Field, utils.MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]), filter.Value)
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space("AND")
			}
		}
	}

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	var total int
	err = service.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		log.Printf("%s", err)
		return nil, 0, 500, utils.ErrInternalServer
	}

	return usersResponse, total, 200, nil
}

func (service *SettingUserService) User(id string) (models.SettingUserResponse, int, error) {
	query := `
	WITH "limited_users" AS (
		SELECT 
			id, name, email, typ, setting_role_id
		FROM "setting.user"
		WHERE id = ?
	)
	SELECT 
		"limited_users".id, 
		"limited_users".name, 
		"limited_users".email, 
		"limited_users".typ,
		COALESCE("setting.role".id, 0), 
		COALESCE("setting.role".name, ''), 
		COALESCE("setting.role".description, ''), 
		COALESCE("setting.permission".id, 0), 
		COALESCE("setting.permission".name, '')
	FROM "limited_users"
	LEFT JOIN "setting.role" ON "limited_users".setting_role_id = "setting.role".id
	LEFT JOIN "setting.role_permission" ON "setting.role".id = "setting.role_permission".setting_role_id
	LEFT JOIN "setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id`

	query, params, err := bqb.New(query, id).ToPgsql()
	if err != nil {
		log.Printf("%s", err)
		return models.SettingUserResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		return models.SettingUserResponse{}, 500, utils.ErrInternalServer
	}

	usersResponse := make([]models.SettingUserResponse, 0)
	permission := make([]models.SettingPermission, 0)
	var lastUser models.SettingUser
	var lastRole models.SettingRole
	for rows.Next() {
		var tmpUser models.SettingUser
		var tmpRole models.SettingRole
		var tmpPermission models.SettingPermission
		err := rows.Scan(&tmpUser.Id, &tmpUser.Name, &tmpUser.Email, &tmpUser.Typ, &tmpRole.Id, &tmpRole.Name, &tmpRole.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			log.Printf("%s", err)
			return models.SettingUserResponse{}, 500, utils.ErrInternalServer
		}

		if lastUser.Id != tmpUser.Id && lastUser.Id != 0 {
			usersResponse = append(usersResponse, models.SettingUserToResponse(lastUser, lastRole, permission))
			lastUser = tmpUser
			lastRole = tmpRole
			permission = make([]models.SettingPermission, 0)
			permission = append(permission, tmpPermission)
			continue
		}

		if lastUser.Id == 0 {
			lastUser = tmpUser
			lastRole = tmpRole
		}

		if tmpPermission.Id != 0 {
			permission = append(permission, tmpPermission)
		}
	}
	if lastUser.Id != 0 {
		usersResponse = append(usersResponse, models.SettingUserToResponse(lastUser, lastRole, permission))
	}

	if len(usersResponse) == 0 {
		return models.SettingUserResponse{}, 404, ErrUserNotFound
	}

	return usersResponse[0], 200, nil
}
