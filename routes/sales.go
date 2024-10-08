package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"system.buon18.com/m/controllers"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/services"
	salesServices "system.buon18.com/m/services/sales"
	"system.buon18.com/m/utils"
)

func SalesRoutes(r *mux.Router, con *database.Connection) {
	controller := controllers.SalesHandler{
		DB: con.DB,
		ServiceFacade: &services.ServiceFacade{
			SalesOrderService:     &salesServices.SalesOrderService{DB: con.DB},
			SalesQuotationService: &salesServices.SalesQuotationService{DB: con.DB},
		},
	}

	r.Handle(
		"/quotations",
		middlewares.Authorize(
			http.HandlerFunc(controller.Quotations),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/quotations/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.Quotation),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/quotations",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateQuotation),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/quotations/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateQuotation),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/quotations/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.DeleteQuotation),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_QUOTATIONS.DELETE},
		),
	).Methods("DELETE", "OPTIONS")

	r.Handle(
		"/orders",
		middlewares.Authorize(
			http.HandlerFunc(controller.Orders),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/orders/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.Order),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/orders",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateOrder),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/orders/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateOrder),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SALES, utils.PREDEFINED_PERMISSIONS.SALES_ORDERS.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")
}
