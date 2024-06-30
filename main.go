package main

import (
	"database/sql"
	"server/config"
	"server/database"
	"server/routes"

	"github.com/gin-gonic/gin"
)

var DB *sql.DB

func main() {
	// -- Initialize
	config.GetConfigInstance()
	DB = database.InitSQL()
	defer DB.Close()

	router := gin.Default()

	routes.Auth(router, DB)
	routes.User(router, DB)
	routes.Account(router, DB)
	routes.Client(router, DB)
	routes.Quote(router, DB)

	router.Routes()
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
