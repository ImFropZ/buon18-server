package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Quote(e *gin.Engine, db *sql.DB) {
	handler := controllers.QuoteHandler{DB: db}

	e.GET("/quotes", handler.List)
	e.GET("/quotes/:id", handler.First)
	e.POST("/quotes", middlewares.Authorize(middlewares.Editor), handler.Create)
	e.POST("/quotes/:id/status", middlewares.Authorize(middlewares.Editor), handler.UpdateStatus)
	e.POST("/quotes/:id/sales-order", middlewares.Authorize(middlewares.Editor), handler.CreateSalesOrder)
	e.PATCH("/quotes/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
	e.DELETE("/quotes/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
	e.DELETE("/quotes/:id/items/:qid", middlewares.Authorize(middlewares.Admin), handler.DeleteItem)
}
