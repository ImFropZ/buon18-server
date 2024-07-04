package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Account(e *gin.Engine, db *sql.DB) {
	handler := controllers.AccountHandler{DB: db}

	e.GET("/accounts", middlewares.Authenticate(db), handler.List)
	e.GET("/accounts/:id", middlewares.Authenticate(db), handler.First)
	e.POST("/accounts", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Editor), handler.Create)
	e.PATCH("/accounts/:id", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Editor), handler.Update)
	e.DELETE("/accounts/:id", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Admin), handler.Delete)
	e.DELETE("/accounts/:id/social-medias/:smid", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Admin), handler.DeleteSocialMedia)
}
