package main

import (
	"server/config"
	"server/database"
	"server/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	database.ConnectDB()
}

func main() {
	router := gin.Default()

	routes.Auth(router, database.DB)

	router.Routes()
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
