package controllers

import (
	"errors"
	"log"
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

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type AuthHandler struct {
	DB *gorm.DB
}

func (handler *AuthHandler) Me(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		log.Printf("Error getting email from context: %v\n", errors.New("email not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get user from db
	var user models.User
	result := handler.DB.First(&user, "email = ?", email)

	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "user doesn't existed"})
		return
	}

	c.JSON(200, gin.H{
		"name":  user.Name,
		"email": user.Email,
		"role":  user.Role,
	})
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
	handler.DB.First(&user, "email = ?", req.Email)

	if user.Deleted {
		c.JSON(401, gin.H{"error": "Your account has been deleted"})
		return
	}

	if user.Email != req.Email || (!utils.ComparePwd(req.Password, user.Pwd) && user.Pwd != "") {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// -- Generate token
	token, err := utils.GenerateWebToken(utils.WebTokenClaims{
		Email: user.Email,
		Role:  user.Role,
	})

	if err != nil {
		log.Printf("Error generating web token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(utils.RefreshTokenClaims{
		Email: user.Email,
	})
	if err != nil {
		log.Printf("Error generating refresh token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
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
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain refresh_token field"))
		return
	}

	token, err := utils.RemoveBearer(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(401, utils.NewErrorResponse(401, "invalid token"))
		return
	}

	// -- Validate token
	claims, err := utils.ValidateWebToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			refreshClaims, refreshErr := utils.ValidateRefreshToken(req.RefreshToken)
			if refreshErr != nil {
				c.JSON(401, utils.NewErrorResponse(401, "invalid refresh token"))
				return
			}

			// -- Check email
			var user models.User
			result := handler.DB.First(&user, "email = ?", refreshClaims.Email)

			if result.RowsAffected == 0 {
				c.JSON(401, utils.NewErrorResponse(401, "email doesn't existed"))
				return
			}

			// -- Generate new token
			token, err := utils.GenerateWebToken(utils.WebTokenClaims{
				Email: user.Email,
				Role:  user.Role,
			})
			if err != nil {
				log.Printf("Error generating web token: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Response new token
			c.JSON(200, utils.NewResponse(200, "success", gin.H{
				"token": token,
			}))
			return
		}

		c.JSON(401, utils.NewErrorResponse(401, "invalid token"))
		return
	}

	token, err = utils.GenerateWebToken(claims)
	if err != nil {
		log.Printf("Error generating web token: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"token": token,
	}))
}

func (handler *AuthHandler) UpdatePassword(c *gin.Context) {
	// -- Parse request
	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain old_password and new_password fields"))
		return
	}

	// -- Get email
	email, _ := c.Get("email")

	// -- Get user from db
	var user models.User
	result := handler.DB.First(&user, "email = ?", email)
	if result.RowsAffected == 0 {
		c.JSON(404, utils.NewErrorResponse(404, "user doesn't existed"))
		return
	}

	// -- Update Pwd if user in db doesn't have pwd
	if user.Pwd != "" {
		// -- Compare pwd
		if ok := utils.ComparePwd(req.OldPassword, user.Pwd); !ok {
			c.JSON(400, utils.NewErrorResponse(400, "invalid old password"))
			return
		}
	}

	// -- Update pwd
	if hashedPwd, err := utils.HashPwd(req.NewPassword); err == nil {
		user.Pwd = hashedPwd
		handler.DB.Save(&user)
		c.JSON(200, utils.NewResponse(200, "success", nil))
	} else {
		log.Printf("Error hashing password: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
	}
}
