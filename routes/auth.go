package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"
	"server/services"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Auth(e *gin.Engine, db *sql.DB) {
	handler := controllers.AuthHandler{
		DB:          db,
		AuthService: &services.AuthService{DB: db},
	}
	authPermissions := utils.PREDEFINED_PERMISSIONS.AUTH

	e.GET("/api/auth/me", middlewares.Authenticate(db), middlewares.Authorize([]string{authPermissions.VIEW}), handler.Me)
	e.POST("/api/auth/login", handler.Login)
	e.POST("/api/auth/refresh-token", handler.RefreshToken)
	e.POST("/api/auth/update-password", middlewares.Authenticate(db), middlewares.Authorize([]string{authPermissions.UPDATE}), handler.UpdatePassword)
	e.PATCH("/api/auth/me", middlewares.Authenticate(db), middlewares.Authorize([]string{authPermissions.UPDATE}), handler.UpdateProfile)
}
