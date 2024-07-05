package main

import (
	"database/sql"
	"server/config"
	"server/database"
	"server/middlewares"
	"server/routes"

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

	rg := router.Group("/api")

	rg.Use(cors.New(cors.Config{
		AllowOrigins:  config.ALLOW_ORIGINS,
		AllowMethods:  config.ALLOW_METHODS,
		AllowHeaders:  config.ALLOW_HEADERS,
		ExposeHeaders: config.EXPOSE_HEADERS,
		MaxAge:        config.MAX_AGE,
	}))

	// -- Routes
	// -- Public
	routes.Auth(rg, DB)

	// -- Private
	rg.Use(middlewares.Authenticate(DB))
	routes.User(rg, DB)
	routes.Account(rg, DB)
	routes.Client(rg, DB)
	routes.Quote(rg, DB)
	routes.SalesOrder(rg, DB)

	router.Routes()

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
