package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func User(e *gin.Engine, db *sql.DB) {
	handler := controllers.UserHandler{DB: db}

	e.GET("/users", middlewares.Authorize(middlewares.Editor), handler.List)
	e.GET("/users/:id", middlewares.Authorize(middlewares.Editor), handler.First)
	e.POST("/users", middlewares.Authorize(middlewares.Admin), handler.Create)
	e.PATCH("/users/:id", middlewares.Authorize(middlewares.Admin), handler.Update)
	e.DELETE("/users/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
}
