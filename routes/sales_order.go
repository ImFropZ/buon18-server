package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func SalesOrder(c *gin.Engine, db *sql.DB) {
	handler := controllers.SalesOrderHandler{DB: db}

	router := c.Group("/sales-orders", middlewares.Authenticate(db))
	{
		router.POST("/:id/status", middlewares.Authorize(middlewares.Editor), handler.UpdateStatus)
	}
}
