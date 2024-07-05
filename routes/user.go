package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func User(e *gin.Engine, db *sql.DB) {
	handler := controllers.UserHandler{DB: db}

	e.GET("/api/users", middlewares.Authorize(middlewares.Editor), handler.List)
	e.GET("/api/users/:id", middlewares.Authorize(middlewares.Editor), handler.First)
	e.POST("/api/users", middlewares.Authorize(middlewares.Admin), handler.Create)
	e.PATCH("/api/users/:id", middlewares.Authorize(middlewares.Admin), handler.Update)
	e.DELETE("/api/users/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
}
