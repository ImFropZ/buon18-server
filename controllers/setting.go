package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"server/models/setting"
	"server/services"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	DB            *sql.DB
	ServiceFacade *services.ServiceFacade
}

func (handler *SettingHandler) Users(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, setting.SettingUserAllowFilterFieldsAndOps, `"setting.user"`).
		PrepareSorts(c, setting.SettingUserAllowSortFields, `"limited_users"`).
		PreparePagination(c)

	users, total, statusCode, err := handler.ServiceFacade.SettingUserService.Users(qp)
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

	user, statusCode, err := handler.ServiceFacade.SettingUserService.User(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"user": user,
	}))
}

func (handler *SettingHandler) CreateUser(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	var user setting.SettingUserCreateRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, err.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(user); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.SettingUserService.CreateUser(&ctx, &user)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "user created successfully", nil))
}

func (handler *SettingHandler) Customers(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, setting.SettingCustomerAllowFilterFieldsAndOps, `"setting.customer"`).
		PrepareSorts(c, setting.SettingCustomerAllowSortFields, `"limited_customers"`).
		PreparePagination(c)

	customers, total, statusCode, err := handler.ServiceFacade.SettingCustomerService.Customers(qp)
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

	customer, statusCode, err := handler.ServiceFacade.SettingCustomerService.Customer(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"customer": customer,
	}))
}

func (handler *SettingHandler) CreateCustomer(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	var customer setting.SettingCustomerCreateRequest
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, err.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(customer); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.SettingCustomerService.CreateCustomer(&ctx, &customer)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "customer created successfully", nil))
}

func (handler *SettingHandler) UpdateCustomer(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	id := c.Param("id")

	var customer setting.SettingCustomerUpdateRequest
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, err.Error()))
		return
	}

	if utils.IsAllFieldsNil(&customer) {
		c.JSON(400, utils.NewErrorResponse(400, "no fields to update"))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(customer); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.SettingCustomerService.UpdateCustomer(&ctx, id, &customer)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "customer updated successfully", nil))
}

func (handler *SettingHandler) Roles(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, setting.SettingRoleAllowFilterFieldsAndOps, `"setting.role"`).
		PrepareSorts(c, setting.SettingRoleAllowSortFields, `"limited_roles"`).
		PreparePagination(c)

	roles, total, statusCode, err := handler.ServiceFacade.SettingRoleService.Roles(qp)
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

	role, statusCode, err := handler.ServiceFacade.SettingRoleService.Role(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"role": role,
	}))
}

func (handler *SettingHandler) CreateRole(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	var role setting.SettingRoleCreateRequest
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, err.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(role); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.SettingRoleService.CreateRole(&ctx, &role)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "role created successfully", nil))
}
