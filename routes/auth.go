package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Auth(rg *gin.RouterGroup, db *sql.DB) {
	handler := controllers.AuthHandler{DB: db}

	rg.GET("/me", middlewares.Authenticate(db), handler.Me)
	rg.POST("/login", handler.Login)
	rg.POST("/refresh", handler.RefreshToken)
	rg.POST("/update-password", middlewares.Authenticate(db), handler.UpdatePassword)
}
