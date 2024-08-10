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
	ErrUserNotFound = errors.New("user not found")
)

type SettingUserService struct {
	DB *sql.DB
}

func (service *SettingUserService) Users(qp *utils.QueryParams) ([]setting.SettingUserResponse, int, int, error) {
	bqbQuery := bqb.New(`
	WITH "limited_users" AS (
		SELECT 
			id, name, email, typ, setting_role_id
		FROM "setting.user"`)

	qp.FilterIntoBqb(bqbQuery)
	qp.PaginationIntoBqb(bqbQuery)

	bqbQuery.Space(`)
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
	LEFT JOIN "setting.permission" ON "setting.role_permission".setting_permission_id = "setting.permission".id`)

	qp.OrderByIntoBqb(bqbQuery, `"limited_users".id ASC, "setting.role".id ASC, "setting.permission".id ASC`)

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

	usersResponse := make([]setting.SettingUserResponse, 0)
	permissions := make([]setting.SettingPermission, 0)
	var lastUser setting.SettingUser
	var lastRole setting.SettingRole
	for rows.Next() {
		var tmpUser setting.SettingUser
		var tmpRole setting.SettingRole
		var tmpPermission setting.SettingPermission
		err := rows.Scan(&tmpUser.Id, &tmpUser.Name, &tmpUser.Email, &tmpUser.Typ, &tmpRole.Id, &tmpRole.Name, &tmpRole.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			log.Printf("%s", err)
			return nil, 0, 500, utils.ErrInternalServer
		}

		if lastUser.Id != tmpUser.Id && lastUser.Id != 0 {
			permissionsResponse := make([]setting.SettingPermissionResponse, 0)
			for _, permission := range permissions {
				permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
			}
			roleResponse := setting.SettingRoleToResponse(lastRole, permissionsResponse)
			usersResponse = append(usersResponse, setting.SettingUserToResponse(lastUser, roleResponse))
			lastUser = tmpUser
			lastRole = tmpRole
			permissions = make([]setting.SettingPermission, 0)
			permissions = append(permissions, tmpPermission)
			continue
		}

		if lastUser.Id == 0 {
			lastUser = tmpUser
			lastRole = tmpRole
		}

		if tmpPermission.Id != 0 {
			permissions = append(permissions, tmpPermission)
		}
	}
	if lastUser.Id != 0 {
		permissionsResponse := make([]setting.SettingPermissionResponse, 0)
		for _, permission := range permissions {
			permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
		}
		roleResponse := setting.SettingRoleToResponse(lastRole, permissionsResponse)
		usersResponse = append(usersResponse, setting.SettingUserToResponse(lastUser, roleResponse))
	}

	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "setting.user"`)
	qp.FilterIntoBqb(bqbQuery)

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

func (service *SettingUserService) User(id string) (setting.SettingUserResponse, int, error) {
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
		return setting.SettingUserResponse{}, 500, utils.ErrInternalServer
	}

	rows, err := service.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		return setting.SettingUserResponse{}, 500, utils.ErrInternalServer
	}

	var user setting.SettingUser
	var role setting.SettingRole
	permissions := make([]setting.SettingPermission, 0)
	for rows.Next() {
		var tmpPermission setting.SettingPermission
		err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Typ, &role.Id, &role.Name, &role.Description, &tmpPermission.Id, &tmpPermission.Name)
		if err != nil {
			log.Printf("%s", err)
			return setting.SettingUserResponse{}, 500, utils.ErrInternalServer
		}

		permissions = append(permissions, tmpPermission)
	}

	if user.Id == 0 {
		return setting.SettingUserResponse{}, 404, ErrUserNotFound
	}

	permissionsResponse := make([]setting.SettingPermissionResponse, 0)
	for _, permission := range permissions {
		permissionsResponse = append(permissionsResponse, setting.SettingPermissionToResponse(permission))
	}
	roleResponse := setting.SettingRoleToResponse(role, permissionsResponse)

	return setting.SettingUserToResponse(user, roleResponse), 200, nil
}
