package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Quote(rg *gin.RouterGroup, db *sql.DB) {
	handler := controllers.QuoteHandler{DB: db}

	rg.GET("/quotes", handler.List)
	rg.GET("/quotes/:id", handler.First)
	rg.POST("/quotes", middlewares.Authorize(middlewares.Editor), handler.Create)
	rg.POST("/quotes/:id/status", middlewares.Authorize(middlewares.Editor), handler.UpdateStatus)
	rg.POST("/quotes/:id/sales-order", middlewares.Authorize(middlewares.Editor), handler.CreateSalesOrder)
	rg.PATCH("/quotes/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
	rg.DELETE("/quotes/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
	rg.DELETE("/quotes/:id/items/:qid", middlewares.Authorize(middlewares.Admin), handler.DeleteItem)
}
