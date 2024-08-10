package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"
	"server/services/sales"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Sales(e *gin.Engine, db *sql.DB) {
	handler := controllers.SalesHandler{
		DB:                    db,
		SalesQuotationService: &sales.SalesQuotationService{DB: db},
		SalesOrderService:     &sales.SalesOrderService{DB: db},
	}

	e.GET("/api/sales/quotations", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.VIEW}), handler.Quotations)
	e.GET("/api/sales/quotations/:id", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.VIEW}), handler.Quotation)
	e.GET("/api/sales/orders", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.VIEW}), handler.Orders)
}
