package controllers

import (
	"database/sql"
	"server/services"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	DB          *sql.DB
	AuthService *services.AuthService
}

func (handler *AuthHandler) Me(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	user, statusCode, err := handler.AuthService.Me(&ctx)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"user": user,
	}))
}

func (handler *AuthHandler) Login(c *gin.Context) {
	// -- Parse request
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain email and password fields"))
		return
	}

	tokenAndRefreshToken, statusCode, err := handler.AuthService.Login(&req)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", tokenAndRefreshToken))
}

func (handler *AuthHandler) RefreshToken(c *gin.Context) {
	// -- Parse request
	var req services.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain refresh_token field"))
		return
	}

	tokenAndRefreshToken, statusCode, err := handler.AuthService.RefreshToken(&req)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", tokenAndRefreshToken))
}

func (handler *AuthHandler) UpdatePassword(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Parse request
	var req services.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain old_password and new_password fields"))
		return
	}

	message, statusCode, err := handler.AuthService.UpdatePassword(&ctx, &req)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, message, nil))
}

func (handler *AuthHandler) UpdateProfile(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Parse request
	var updateData services.UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request"))
		return
	}
	if utils.IsAllFieldsNil(&updateData) {
		c.JSON(400, utils.NewErrorResponse(400, "request body should contain at least one field"))
		return
	}

	message, statusCode, err := handler.AuthService.UpdateProfile(&ctx, &updateData)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, message, nil))
}
