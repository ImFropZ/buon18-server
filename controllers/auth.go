package controllers

import (
	"errors"
	"fmt"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type AuthHandler struct {
	DB *gorm.DB
}

func (handler *AuthHandler) Login(c *gin.Context) {
	// -- Parse request
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// -- Validate user
	var user models.User
	handler.DB.First(&user)

	// -- Hash password
	hashPwd, err := utils.HashPwd(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if user.Email != req.Email || (user.Pwd != hashPwd && user.Pwd != "") {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// -- Generate token
	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email: user.Email,
		Role:  "admin",
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	// -- Parse request
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.RemoveBearer(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	// -- Validate token
	claims, err := utils.ValidateWebToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			refreshClaims, refreshErr := utils.ValidateRefreshToken(req.RefreshToken)
			if refreshErr != nil {
				fmt.Printf("Error parsing token: %v\n", refreshErr)
				c.JSON(401, gin.H{"error": refreshErr.Error()})
				return
			}

			// -- Check email
			var user models.User
			result := handler.DB.First(&user, "email = ?", refreshClaims.Email)

			if result.RowsAffected == 0 {
				c.JSON(401, gin.H{"error": "User not found"})
				return
			}

			// -- Generate new token
			token, err := utils.GenerateWebToken(utils.WebTokenClaims{
				Email: user.Email,
				Role:  "admin",
			})
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			// -- Response new token
			c.JSON(200, gin.H{
				"token": token,
			})
			return
		}

		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	token, err = utils.GenerateWebToken(claims)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
