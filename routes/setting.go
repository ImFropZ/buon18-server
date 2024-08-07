package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"
	"server/services/setting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Setting(e *gin.Engine, db *sql.DB) {
	handler := controllers.SettingHandler{
		DB:                     db,
		SettingUserService:     &setting.SettingUserService{DB: db},
		SettingCustomerService: &setting.SettingCustomerService{DB: db},
		SettingRoleService:     &setting.SettingRoleService{DB: db},
	}

	e.GET("/api/setting/users", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.VIEW}), handler.Users)
	e.GET("/api/setting/users/:id", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.VIEW}), handler.User)
	e.GET("/api/setting/customers", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.VIEW}), handler.Customers)
	e.GET("/api/setting/customers/:id", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.VIEW}), handler.Customer)
	e.GET("/api/setting/roles", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.VIEW}), handler.Roles)
	e.GET("/api/setting/roles/:id", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.VIEW}), handler.Role)
}
