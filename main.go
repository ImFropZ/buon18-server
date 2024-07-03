package main

import (
	"database/sql"
	"server/config"
	"server/database"
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

	router.Use(cors.New(cors.Config{
		AllowOrigins:  config.ALLOW_ORIGINS,
		AllowMethods:  config.ALLOW_METHODS,
		AllowHeaders:  config.ALLOW_HEADERS,
		ExposeHeaders: config.EXPOSE_HEADERS,
		MaxAge:        config.MAX_AGE,
	}))
	router.SetTrustedProxies(config.TRUSTED_PROXIES)

	routes.Auth(router, DB)
	routes.User(router, DB)
	routes.Account(router, DB)
	routes.Client(router, DB)
	routes.Quote(router, DB)
	routes.SalesOrder(router, DB)

	router.Routes()

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
