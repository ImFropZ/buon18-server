package routes

import (
	"server/controllers"
	"server/database"
	"server/middlewares"
	"server/models/sales"
	service "server/services/sales"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Sales(e *gin.Engine, connection *database.Connection) {
	handler := controllers.SalesHandler{
		DB:                    connection.DB,
		SalesQuotationService: &service.SalesQuotationService{DB: connection.DB},
		SalesOrderService:     &service.SalesOrderService{DB: connection.DB},
	}

	e.GET(
		"/api/sales/quotations",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.VIEW}),
		middlewares.ValkeyCache[[]sales.SalesQuotationResponse](connection, "quotations"),
		handler.Quotations,
	)
	e.GET(
		"/api/sales/quotations/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.VIEW}),
		handler.Quotation,
	)
	e.GET(
		"/api/sales/orders",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.VIEW}),
		middlewares.ValkeyCache[[]sales.SalesOrderResponse](connection, "orders"),
		handler.Orders,
	)
	e.GET(
		"/api/sales/orders/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.VIEW}),
		handler.Order,
	)
}
