package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func User(e *gin.Engine, db *sql.DB) {
	handler := controllers.UserHandler{DB: db}

	router := e.Group("/users", middlewares.Authenticate(db))
	{
		router.GET("/", middlewares.Authorize(middlewares.Editor), handler.List)
		router.GET("/:id", middlewares.Authorize(middlewares.Editor), handler.First)
		router.POST("/", middlewares.Authorize(middlewares.Admin), handler.Create)
		router.PATCH("/:id", middlewares.Authorize(middlewares.Admin), handler.Update)
		router.DELETE("/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
	}
}
