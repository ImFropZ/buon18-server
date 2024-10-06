package routes

import (
	"system.buon18.com/m/controllers"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/models/sales"
	"system.buon18.com/m/services"
	salesServices "system.buon18.com/m/services/sales"
	"system.buon18.com/m/utils"

	"github.com/gin-gonic/gin"
)

func Sales(e *gin.Engine, connection *database.Connection) {
	handler := controllers.SalesHandler{
		DB: connection.DB,
		ServiceFacade: &services.ServiceFacade{
			SalesOrderService:     &salesServices.SalesOrderService{DB: connection.DB},
			SalesQuotationService: &salesServices.SalesQuotationService{DB: connection.DB},
		},
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
	e.POST(
		"/api/sales/quotations",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.CREATE}),
		handler.CreateQuotation,
	)
	e.PATCH(
		"/api/sales/quotations/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.UPDATE}),
		handler.UpdateQuotation,
	)
	e.DELETE(
		"/api/sales/quotations/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.DELETE}),
		handler.DeleteQuotation,
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
	e.POST(
		"/api/sales/orders",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.CREATE}),
		handler.CreateOrder,
	)
	e.PATCH(
		"/api/sales/orders/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.UPDATE}),
		handler.UpdateOrder,
	)
}
