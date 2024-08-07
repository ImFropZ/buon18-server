package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(e *gin.Engine, db *sql.DB) {
	handler := controllers.AuthHandler{DB: db}

	e.GET("/api/auth/me", middlewares.Authenticate(db), middlewares.Authorize([]string{"VIEW_PROFILE"}), handler.Me)
	e.POST("/api/auth/login", handler.Login)
	e.POST("/api/auth/refresh-token", handler.RefreshToken)
	e.POST("/api/auth/update-password", middlewares.Authenticate(db), middlewares.Authorize([]string{"UPDATE_PROFILE"}), handler.UpdatePassword)
	e.PATCH("/api/auth/me", middlewares.Authenticate(db), middlewares.Authorize([]string{"UPDATE_PROFILE"}), handler.UpdateProfile)
}
