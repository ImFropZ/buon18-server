package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func User(rg *gin.RouterGroup, db *sql.DB) {
	handler := controllers.UserHandler{DB: db}

	rg.GET("/users", middlewares.Authorize(middlewares.Editor), handler.List)
	rg.GET("/users/:id", middlewares.Authorize(middlewares.Editor), handler.First)
	rg.POST("/users", middlewares.Authorize(middlewares.Admin), handler.Create)
	rg.PATCH("/users/:id", middlewares.Authorize(middlewares.Admin), handler.Update)
	rg.DELETE("/users/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
}
