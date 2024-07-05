package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func SalesOrder(rg *gin.RouterGroup, db *sql.DB) {
	handler := controllers.SalesOrderHandler{DB: db}

	rg.GET("/sales-orders", handler.List)
	rg.GET("/sales-orders/:id", handler.First)
	rg.POST("/sales-orders/:id/invoice", middlewares.Authorize(middlewares.Editor), handler.CreateInvoice)
	rg.POST("/sales-orders/:id/status", middlewares.Authorize(middlewares.Editor), handler.UpdateStatus)
}
