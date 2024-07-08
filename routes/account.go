package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Account(e *gin.Engine, db *sql.DB) {
	handler := controllers.AccountHandler{DB: db}

	e.GET("/api/accounts", middlewares.Authenticate(db), handler.List)
	e.GET("/api/accounts/:id", middlewares.Authenticate(db), handler.First)
	e.POST("/api/accounts", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Editor), handler.Create)
	e.PATCH("/api/accounts/:id", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Editor), handler.Update)
	e.DELETE("/api/accounts/:id", middlewares.Authenticate(db), middlewares.Authorize(middlewares.Admin), handler.Delete)
}
