package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Auth(e *gin.Engine, db *gorm.DB) {
	handler := controllers.AuthHandler{DB: db}

	e.POST("/login", handler.Login)
	e.POST("/refresh", handler.RefreshToken)
}
