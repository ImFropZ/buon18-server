package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Quote(e *gin.Engine, db *sql.DB) {
	handler := controllers.QuoteHandler{DB: db}

	e.GET("/api/quotes", handler.List)
	e.GET("/api/quotes/:id", handler.First)
	e.POST("/api/quotes", middlewares.Authorize(middlewares.Editor), handler.Create)
	e.POST("/api/quotes/:id/status", middlewares.Authorize(middlewares.Editor), handler.UpdateStatus)
	e.POST("/api/quotes/:id/sales-order", middlewares.Authorize(middlewares.Editor), handler.CreateSalesOrder)
	e.PATCH("/api/quotes/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
	e.DELETE("/api/quotes/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
}
