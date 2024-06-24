package routes

import (
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Account(e *gin.Engine, db *gorm.DB) {
	handler := controllers.AccountHandler{DB: db}

	account := e.Group("/accounts", middlewares.Authenticate())
	{
		account.GET("/", middlewares.Authorize(middlewares.User), handler.List)
		account.GET("/:id", middlewares.Authorize(middlewares.User), handler.First)
		account.POST("/", middlewares.Authorize(middlewares.Editor), handler.Create)
	}
}
