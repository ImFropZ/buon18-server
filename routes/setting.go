package routes

import (
	"system.buon18.com/m/controllers"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/models/setting"
	"system.buon18.com/m/services"
	settingServices "system.buon18.com/m/services/setting"
	"system.buon18.com/m/utils"

	"github.com/gin-gonic/gin"
)

func Setting(e *gin.Engine, connection *database.Connection) {
	handler := controllers.SettingHandler{
		DB: connection.DB,
		ServiceFacade: &services.ServiceFacade{
			SettingCustomerService:   &settingServices.SettingCustomerService{DB: connection.DB},
			SettingRoleService:       &settingServices.SettingRoleService{DB: connection.DB},
			SettingUserService:       &settingServices.SettingUserService{DB: connection.DB},
			SettingPermissionService: &settingServices.SettingPermissionService{DB: connection.DB},
		},
	}

	e.GET(
		"/api/setting/users",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.VIEW}),
		middlewares.ValkeyCache[[]setting.SettingUserResponse](connection, "users"),
		handler.Users,
	)
	e.GET(
		"/api/setting/users/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.VIEW}),
		handler.User,
	)
	e.POST(
		"/api/setting/users",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.CREATE}),
		handler.CreateUser,
	)
	e.PATCH(
		"/api/setting/users/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.UPDATE}),
		handler.UpdateUser,
	)
	e.GET(
		"/api/setting/customers",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.VIEW}),
		middlewares.ValkeyCache[[]setting.SettingCustomerResponse](connection, "customers"),
		handler.Customers,
	)
	e.GET(
		"/api/setting/customers/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.VIEW}),
		handler.Customer,
	)
	e.POST(
		"/api/setting/customers",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.CREATE}),
		handler.CreateCustomer,
	)
	e.PATCH(
		"/api/setting/customers/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.UPDATE}),
		handler.UpdateCustomer,
	)
	e.DELETE(
		"/api/setting/customers/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.DELETE}),
		handler.DeleteCustomer,
	)
	e.GET(
		"/api/setting/roles",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.VIEW}),
		middlewares.ValkeyCache[[]setting.SettingRoleResponse](connection, "roles"),
		handler.Roles,
	)
	e.GET(
		"/api/setting/roles/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.VIEW}),
		handler.Role,
	)
	e.POST(
		"/api/setting/roles",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.CREATE}),
		handler.CreateRole,
	)
	e.PATCH(
		"/api/setting/roles/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.UPDATE}),
		handler.UpdateRole,
	)
	e.DELETE(
		"/api/setting/roles/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.DELETE}),
		handler.DeleteRole,
	)
	e.GET(
		"/api/setting/permissions",
		handler.Permissions,
	)
}
