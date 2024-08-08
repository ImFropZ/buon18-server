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
	DB                     *sql.DB
	SettingUserService     setting.SettingUserService
	SettingCustomerService setting.SettingCustomerService
	SettingRoleService     setting.SettingRoleService
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

func (handler *SettingHandler) User(c *gin.Context) {
	id := c.Param("id")

	user, statusCode, err := handler.SettingUserService.User(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"user": user,
	}))
}

func (handler *SettingHandler) Customers(c *gin.Context) {
	qp := utils.NewQueryParams()
	for _, filter := range models.SettingCustomerAllowFilterFieldsAndOps {
		if validFilter, ok := c.GetQuery(filter); ok {
			qp.AddFilter(fmt.Sprintf(`"setting.customer".%s=%s`, filter, validFilter))
		}
	}
	for _, sort := range models.SettingCustomerAllowSortFields {
		if validSort, ok := c.GetQuery(fmt.Sprintf("sort-%s", sort)); ok {
			qp.AddOrderBy(fmt.Sprintf(`LOWER("setting.customer".%s) %s`, sort, validSort))
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

	customers, total, statusCode, err := handler.SettingCustomerService.Customers(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"customers": customers,
	}))
}

func (handler *SettingHandler) Customer(c *gin.Context) {
	id := c.Param("id")

	customer, statusCode, err := handler.SettingCustomerService.Customer(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"customer": customer,
	}))
}

func (handler *SettingHandler) Roles(c *gin.Context) {
	qp := utils.NewQueryParams()
	for _, filter := range models.SettingRoleAllowFilterFieldsAndOps {
		if validFilter, ok := c.GetQuery(filter); ok {
			qp.AddFilter(fmt.Sprintf(`"setting.role".%s=%s`, filter, validFilter))
		}
	}
	for _, sort := range models.SettingRoleAllowSortFields {
		if validSort, ok := c.GetQuery(fmt.Sprintf("sort-%s", sort)); ok {
			qp.AddOrderBy(fmt.Sprintf(`LOWER("limited_roles".%s) %s`, sort, validSort))
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

	roles, total, statusCode, err := handler.SettingRoleService.Roles(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"roles": roles,
	}))
}
