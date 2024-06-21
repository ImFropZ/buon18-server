package controllers

import (
	"fmt"
	"server/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

type User struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Pwd   string `json:"password"`
}

func (User) TableName() string {
	return "user"
}

type AuthHandler struct {
	DB *gorm.DB
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user User
	handler.DB.First(&user)

	fmt.Printf("User: %+v\n", user)

	token, err := utils.GenerateWebToken(utils.Claims{
		Email: user.Email,
		Role:  "admin",
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.ValidateWebToken(req.Token)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateWebToken(claims)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
