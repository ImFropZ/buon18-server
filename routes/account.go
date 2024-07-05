package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Account(rg *gin.RouterGroup, db *sql.DB) {
	handler := controllers.AccountHandler{DB: db}

	rg.GET("/accounts", middlewares.Authenticate(db), handler.List)
	rg.GET("/accounts/:id", middlewares.Authenticate(db), handler.First)
	rg.POST("/accounts", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Editor), handler.Create)
	rg.PATCH("/accounts/:id", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Editor), handler.Update)
	rg.DELETE("/accounts/:id", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Admin), handler.Delete)
	rg.DELETE("/accounts/:id/social-medias/:smid", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Admin), handler.DeleteSocialMedia)
}
