package routes

import (
	"server/controllers"
	"server/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Auth(e *gin.Engine, db *gorm.DB) {
	handler := controllers.AuthHandler{DB: db}

	e.GET("/me", middlewares.Authenticate(), handler.Me)
	e.POST("/login", handler.Login)
	e.POST("/refresh", handler.RefreshToken)
	e.POST("/update-password", middlewares.Authenticate(), handler.UpdatePassword)
}
