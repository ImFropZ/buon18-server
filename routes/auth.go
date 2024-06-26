package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(e *gin.Engine, db *sql.DB) {
	handler := controllers.AuthHandler{DB: db}

	e.GET("/me", middlewares.Authenticate(db), handler.Me)
	e.POST("/login", handler.Login)
	e.POST("/refresh", handler.RefreshToken)
	e.POST("/update-password", middlewares.Authenticate(db), handler.UpdatePassword)
}
