package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"system.buon18.com/m/controllers"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"
)

func AuthRoutes(r *mux.Router, con *database.Connection) {
	controller := controllers.AuthHandler{
		DB: con.DB,
		ServiceFacade: &services.ServiceFacade{
			AuthService: &services.AuthService{DB: con.DB},
		},
	}
	authPermissions := utils.PREDEFINED_PERMISSIONS.AUTH

	r.Handle(
		"/me",
		middlewares.Authenticate(
			middlewares.Authorize(
				http.HandlerFunc(controller.Me),
				[]string{authPermissions.VIEW},
			),
			con.DB,
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/login",
		http.HandlerFunc(controller.Login),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/refresh-token",
		http.HandlerFunc(controller.RefreshToken),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/update-password",
		middlewares.Authenticate(
			middlewares.Authorize(
				http.HandlerFunc(controller.UpdatePassword),
				[]string{authPermissions.UPDATE},
			),
			con.DB,
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/me",
		middlewares.Authenticate(
			middlewares.Authorize(
				http.HandlerFunc(controller.UpdateProfile),
				[]string{authPermissions.UPDATE},
			),
			con.DB,
		),
	).Methods("PATCH", "OPTIONS")
}
