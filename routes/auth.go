package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(e *gin.Engine, db *sql.DB) {
	handler := controllers.AuthHandler{DB: db}

	e.GET("/api/me", middlewares.Authenticate(db), handler.Me)
	e.POST("/api/login", handler.Login)
	e.POST("/api/refresh", handler.RefreshToken)
	e.POST("/api/update-password", middlewares.Authenticate(db), handler.UpdatePassword)
}
