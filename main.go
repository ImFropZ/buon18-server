package main

import (
	"database/sql"
	"log"
	"net/http"
	"server/config"
	"server/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var DB *sql.DB

func main() {
	// -- Initialize
	config := config.GetConfigInstance()
	DB = database.InitSQL()
	defer DB.Close()

	router := gin.Default()
	router.SetTrustedProxies(config.TRUSTED_PROXIES)

	router.Use(cors.New(cors.Config{
		AllowOrigins:  config.ALLOW_ORIGINS,
		AllowMethods:  config.ALLOW_METHODS,
		AllowHeaders:  config.ALLOW_HEADERS,
		ExposeHeaders: config.EXPOSE_HEADERS,
		MaxAge:        config.MAX_AGE,
	}))

	// -- Routes
	// -- Public
	// routes.Auth(router, DB)

	// -- Private
	// router.Use(middlewares.Authenticate(DB))
	// routes.User(router, DB)
	// routes.Account(router, DB)
	// routes.Client(router, DB)
	// routes.Quote(router, DB)
	// routes.SalesOrder(router, DB)

	router.Routes()

	// -- Start HTTP Server
	if config.CERT_FILE == "" || config.KEY_FILE == "" {
		server := &http.Server{
			Addr:    ":80",
			Handler: router,
		}

		log.Printf("Server started at port %s\n", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	} else {
		server := &http.Server{
			Addr:    ":443",
			Handler: router,
		}

		log.Printf("Server started at port %s\n", server.Addr)

		if err := server.ListenAndServeTLS(config.CERT_FILE, config.KEY_FILE); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}
}
