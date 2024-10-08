package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"system.buon18.com/m/controllers"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/services"
	settingServices "system.buon18.com/m/services/setting"
	"system.buon18.com/m/utils"
)

func SettingRoutes(r *mux.Router, con *database.Connection) {
	controller := controllers.SettingHandler{
		DB: con.DB,
		ServiceFacade: &services.ServiceFacade{
			SettingCustomerService:   &settingServices.SettingCustomerService{DB: con.DB},
			SettingRoleService:       &settingServices.SettingRoleService{DB: con.DB},
			SettingUserService:       &settingServices.SettingUserService{DB: con.DB},
			SettingPermissionService: &settingServices.SettingPermissionService{DB: con.DB},
		},
	}

	r.Handle(
		"/users",
		middlewares.Authorize(
			http.HandlerFunc(controller.Users),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/users/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.User),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/users",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateUser),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/users/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateUser),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_USERS.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/customers",
		middlewares.Authorize(
			http.HandlerFunc(controller.Customers),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/customers/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.Customer),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/customers",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateCustomer),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/customers/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateCustomer),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/customers/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.DeleteCustomer),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_CUSTOMERS.DELETE},
		),
	).Methods("DELETE", "OPTIONS")

	r.Handle(
		"/roles",
		middlewares.Authorize(
			http.HandlerFunc(controller.Roles),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/roles/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.Role),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/roles",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateRole),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/roles/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateRole),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/roles/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.DeleteRole),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_SETTING, utils.PREDEFINED_PERMISSIONS.SETTING_ROLES.DELETE},
		),
	).Methods("DELETE", "OPTIONS")

	r.Handle(
		"/permissions",
		http.HandlerFunc(controller.Permissions),
	).Methods("GET", "OPTIONS")
}
