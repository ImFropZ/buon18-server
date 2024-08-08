package controllers

import (
	"database/sql"
	"fmt"
	"server/models"
	"server/services/setting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	DB                 *sql.DB
	SettingUserService setting.SettingUserService
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

	users, total, statusCode, err := handler.SettingUserService.Users(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"users": users,
	}))
}
