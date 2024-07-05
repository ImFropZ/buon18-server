package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Client(e *gin.Engine, db *sql.DB) {
	handler := controllers.ClientHandler{DB: db}

	e.GET("/api/clients", handler.List)
	e.GET("/api/clients/:id", handler.First)
	e.POST("/api/clients", middlewares.Authorize(middlewares.Editor), handler.Create)
	e.PATCH("/api/clients/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
	e.DELETE("/api/clients/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
	e.DELETE("/api/clients/:id/social-medias/:smid", middlewares.Authorize(middlewares.Admin), handler.DeleteSocialMedia)
}
