package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
)

func Client(e *gin.Engine, db *sql.DB) {
	handler := controllers.ClientHandler{DB: db}

	account := e.Group("/clients", middlewares.Authenticate(db))
	{
		account.GET("/:id", middlewares.Authorize(middlewares.User), handler.First)
	}
}
