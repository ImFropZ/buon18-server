package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
)

type SettingHandler struct {
	DB *sql.DB
}

func (handler *SettingHandler) Users(c *gin.Context) {
	qp := utils.NewQueryParams()
	for _, filter := range models.SettingUserAllowFilterFieldsAndOps {
		if validFilter, ok := c.GetQuery(filter); ok {
			qp.AddFilter(fmt.Sprintf(`"setting.user".%s=%s`, filter, validFilter))
		}
	}
	for _, sort := range models.SettingUserAllowSortFields {
		if validSort, ok := c.GetQuery(fmt.Sprintf("sort-%s", sort)); ok {
			qp.AddOrderBy(fmt.Sprintf(`LOWER("limited_users".%s) %s`, sort, validSort))
		}
	}
	for _, pagination := range []string{"offset", "limit"} {
		if validPagination, ok := c.GetQuery(pagination); ok {
			if pagination == "offset" {
				qp.AddOffset(utils.StrToInt(validPagination, 0))
			} else {
				qp.AddLimit(utils.StrToInt(validPagination, 10))
			}
		}
	}

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
		c.JSON(500, utils.NewErrorResponse(500, "Internal Server Error"))
		return
	}

	rows, err := handler.DB.Query(query, params...)
	if err != nil {
		log.Printf("%s", err)
		c.JSON(500, utils.NewErrorResponse(500, "Internal Server Error"))
		return
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
			c.JSON(500, utils.NewErrorResponse(500, "Internal Server Error"))
			return
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
		c.JSON(500, utils.NewErrorResponse(500, "Internal Server Error"))
		return
	}

	var total int
	err = handler.DB.QueryRow(query, params...).Scan(&total)
	if err != nil {
		log.Printf("%s", err)
		c.JSON(500, utils.NewErrorResponse(500, "Internal Server Error"))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"users": usersResponse,
	}))
}
