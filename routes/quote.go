package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Quote(e *gin.Engine, db *sql.DB) {
	handler := controllers.QuoteHandler{DB: db}

	router := e.Group("/quotes", middlewares.Authenticate(db))
	{
		router.GET("/", handler.List)
		router.GET("/:id", handler.First)
		router.POST("/", middlewares.Authorize(middlewares.Editor), handler.Create)
		router.POST("/:id/status", middlewares.Authorize(middlewares.Editor), handler.UpdateStatus)
		router.POST("/:id/sales-order", middlewares.Authorize(middlewares.Editor), handler.CreateSalesOrder)
		router.PATCH("/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
		router.DELETE("/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
		router.DELETE("/:id/items/:qid", middlewares.Authorize(middlewares.Admin), handler.DeleteItem)
	}
}
