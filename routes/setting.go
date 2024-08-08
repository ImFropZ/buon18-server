package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Setting(e *gin.Engine, db *sql.DB) {
	handler := controllers.SettingHandler{DB: db}

	e.GET("/api/setting/users", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.VIEW}), handler.Users)
}
