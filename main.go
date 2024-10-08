package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"system.buon18.com/m/config"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/routes"

	"github.com/gorilla/mux"
)

var DB *sql.DB

func main() {
	// -- Initialize
	config := config.GetConfigInstance()
	DB = database.InitSQL(config.DB_CONNECTION_STRING)
	defer DB.Close()

	// -- Initialize Valkey
	valkeyClient := database.InitValkey(config.VALKEY_ADDRESSES, config.VALKEY_PWD)
	if valkeyClient != nil {
		defer (*valkeyClient).Close()
	}

	connection := database.Connection{
		DB:     DB,
		Valkey: valkeyClient,
	}

	r := mux.NewRouter()

	r.Use(
		middlewares.CORSHandler,
		middlewares.LoggerHandler,
		middlewares.ErrorHandler,
	)

	apiRoute := r.PathPrefix("/api").Subrouter()

	routes.AuthRoutes(apiRoute.PathPrefix("/auth").Subrouter(), &connection)

	// Auth Routes
	apiRoute.Use(func(next http.Handler) http.Handler {
		return middlewares.Authenticate(next, DB)
	})
	routes.SettingRoutes(apiRoute.PathPrefix("/setting").Subrouter(), &connection)
	routes.SalesRoutes(apiRoute.PathPrefix("/sales").Subrouter(), &connection)
	routes.AccountingRoutes(apiRoute.PathPrefix("/accounting").Subrouter(), &connection)

	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		slog.Info(fmt.Sprintf("%s\t%s\n", methods, t))
		return nil
	})

	// -- Start HTTP Server
	if config.CERT_FILE == "" || config.KEY_FILE == "" {
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", config.PORT),
			Handler: r,
		}

		log.Printf("Server started at port %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	} else {
		server := &http.Server{
			Addr:    ":443",
			Handler: r,
		}

		log.Printf("Server started at port %s\n", server.Addr)

		if err := server.ListenAndServeTLS(config.CERT_FILE, config.KEY_FILE); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}
}
