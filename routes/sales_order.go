package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func SalesOrder(e *gin.Engine, db *sql.DB) {
	handler := controllers.SalesOrderHandler{DB: db}

	e.GET("/api/sales-orders", handler.List)
	e.GET("/api/sales-orders/:id", handler.First)
	e.POST("/api/sales-orders/:id/invoice", middlewares.Authorize(middlewares.Editor), handler.CreateInvoice)
	e.POST("/api/sales-orders/:id/status", middlewares.Authorize(middlewares.Editor), handler.UpdateStatus)
}
