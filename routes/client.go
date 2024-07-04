package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Client(e *gin.Engine, db *sql.DB) {
	handler := controllers.ClientHandler{DB: db}

	e.GET("/clients", handler.List)
	e.GET("/clients/:id", handler.First)
	e.POST("/clients", middlewares.Authorize(middlewares.Editor), handler.Create)
	e.PATCH("/clients/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
	e.DELETE("/clients/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
	e.DELETE("/clients/:id/social-medias/:smid", middlewares.Authorize(middlewares.Admin), handler.DeleteSocialMedia)
}
