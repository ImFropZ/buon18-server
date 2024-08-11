package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"server/models/setting"
	services "server/services/setting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	DB                     *sql.DB
	SettingUserService     *services.SettingUserService
	SettingCustomerService *services.SettingCustomerService
	SettingRoleService     *services.SettingRoleService
}

func (handler *SettingHandler) Users(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, setting.SettingUserAllowFilterFieldsAndOps, `"setting.user"`).
		PrepareSorts(c, setting.SettingUserAllowSortFields, `"limited_users"`).
		PreparePagination(c)

	users, total, statusCode, err := handler.SettingUserService.Users(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"users": users,
	}))

	c.Set("total", total)
	if usersByte, err := json.Marshal(users); err == nil {
		c.Set("response", usersByte)
	}
}

func (handler *SettingHandler) User(c *gin.Context) {
	id := c.Param("id")

	user, statusCode, err := handler.SettingUserService.User(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"user": user,
	}))
}

func (handler *SettingHandler) Customers(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, setting.SettingCustomerAllowFilterFieldsAndOps, `"setting.customer"`).
		PrepareSorts(c, setting.SettingCustomerAllowSortFields, `"limited_customers"`).
		PreparePagination(c)

	customers, total, statusCode, err := handler.SettingCustomerService.Customers(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"customers": customers,
	}))

	c.Set("total", total)
	if customersByte, err := json.Marshal(customers); err == nil {
		c.Set("response", customersByte)
	}
}

func (handler *SettingHandler) Customer(c *gin.Context) {
	id := c.Param("id")

	customer, statusCode, err := handler.SettingCustomerService.Customer(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"customer": customer,
	}))
}

func (handler *SettingHandler) Roles(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, setting.SettingRoleAllowFilterFieldsAndOps, `"setting.role"`).
		PrepareSorts(c, setting.SettingRoleAllowSortFields, `"limited_roles"`).
		PreparePagination(c)

	roles, total, statusCode, err := handler.SettingRoleService.Roles(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"roles": roles,
	}))

	c.Set("total", total)
	if rolesByte, err := json.Marshal(roles); err == nil {
		c.Set("response", rolesByte)
	}
}

func (handler *SettingHandler) Role(c *gin.Context) {
	id := c.Param("id")

	role, statusCode, err := handler.SettingRoleService.Role(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"role": role,
	}))
}
