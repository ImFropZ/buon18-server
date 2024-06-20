package main

import (
	"net/http"

	"server/initializers"
	"server/routes"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	routes.Account(router)

	router.Routes()
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
