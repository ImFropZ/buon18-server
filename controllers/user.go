package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func (handler *UserHandler) List(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "List users",
	})
}
