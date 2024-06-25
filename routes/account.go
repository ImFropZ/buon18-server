package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Account(e *gin.Engine, db *sql.DB) {
	handler := controllers.AccountHandler{DB: db}

	account := e.Group("/accounts", middlewares.Authenticate(db))
	{
		account.GET("/", middlewares.Authorize(middlewares.User), handler.List)
		account.GET("/:id", middlewares.Authorize(middlewares.User), handler.First)
		account.POST("/", middlewares.Authorize(middlewares.Editor), handler.Create)
	}
}
