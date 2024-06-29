package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Client(e *gin.Engine, db *sql.DB) {
	handler := controllers.ClientHandler{DB: db}

	client := e.Group("/clients", middlewares.Authenticate(db))
	{
		client.GET("/", middlewares.Authorize(middlewares.User), handler.List)
		client.GET("/:id", middlewares.Authorize(middlewares.User), handler.First)
		client.POST("/", middlewares.Authorize(middlewares.Editor), handler.Create)
		client.PATCH("/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
		client.DELETE("/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
		client.DELETE("/:id/social-medias/:smid", middlewares.Authorize(middlewares.Admin), handler.DeleteSocialMedia)
	}
}
