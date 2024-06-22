package controllers

import (
	"server/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaginationQueryParams struct {
	Offset int
	Limit  int
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UserHandler struct {
	DB *gorm.DB
}

func (handler *UserHandler) List(c *gin.Context) {
	queryParams := PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	if offset, err := strconv.Atoi(c.Query("offset")); err == nil {
		if offset < 0 {
			offset = 0
		}
		queryParams.Offset = offset
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		if limit < 1 {
			limit = 10
		}
		queryParams.Limit = limit
	}

	var users []UserResponse
	result := handler.DB.Model(&models.User{}).Limit(queryParams.Limit).Offset(queryParams.Offset).Find(&users)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
	})
}
