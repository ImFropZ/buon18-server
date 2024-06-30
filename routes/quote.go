package routes

import (
	"database/sql"
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func Quote(e *gin.Engine, db *sql.DB) {
	handler := controllers.QuoteHandler{DB: db}

	router := e.Group("/quotes")
	{
		router.GET("/", handler.List)
	}
}
