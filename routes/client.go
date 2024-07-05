package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Client(rg *gin.RouterGroup, db *sql.DB) {
	handler := controllers.ClientHandler{DB: db}

	rg.GET("/clients", handler.List)
	rg.GET("/clients/:id", handler.First)
	rg.POST("/clients", middlewares.Authorize(middlewares.Editor), handler.Create)
	rg.PATCH("/clients/:id", middlewares.Authorize(middlewares.Editor), handler.Update)
	rg.DELETE("/clients/:id", middlewares.Authorize(middlewares.Admin), handler.Delete)
	rg.DELETE("/clients/:id/social-medias/:smid", middlewares.Authorize(middlewares.Admin), handler.DeleteSocialMedia)
}
