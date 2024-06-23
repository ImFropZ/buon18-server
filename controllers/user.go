package controllers

import (
	"errors"
	"fmt"
	"server/models"
	"server/utils"
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

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
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

func (handler *UserHandler) Create(c *gin.Context) {
	// -- Get email
	email, _ := c.Get("email")
	if email == nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// -- Bind request
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// -- Validate role
	role, ok := utils.ValidateRole(request.Role)
	if !ok {
		c.JSON(400, gin.H{
			"error": "invalid role",
		})
		return
	}

	// -- Hash password
	hashedPwd, err := utils.HashPwd(request.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// -- Prepare user model
	user := models.User{
		Name:  request.Name,
		Email: request.Email,
		Pwd:   hashedPwd,
		Role:  role,
	}

	// -- Query user by email
	var existingUser models.User
	result := handler.DB.Where("email = ?", email).First(&existingUser)
	if result.Error != nil {
		fmt.Printf("existingUser: %v\n", result.Error)
		c.JSON(400, gin.H{
			"error": "email already exists",
		})
		return
	}

	// -- Prepare for create
	if err := user.PrepareForCreate(existingUser.ID, existingUser.ID); err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// -- Create user
	result = handler.DB.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(400, gin.H{
				"error": "email already exists",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(201, gin.H{
		"user": UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	})
}

func (handler *UserHandler) Delete(c *gin.Context) {
	// -- Get email
	email, _ := c.Get("email")
	if email == nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// -- Get user ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	// -- Query user by ID
	var user models.User
	result := handler.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		c.JSON(404, gin.H{
			"error": "user not found",
		})
		return
	}

	// -- Check if user is deleting itself
	if user.Email == email {
		c.JSON(400, gin.H{
			"error": "cannot delete yourself",
		})
		return
	}

	// -- Check if user already deleted
	if user.Deleted {
		c.JSON(400, gin.H{
			"error": "user already deleted",
		})
		return
	}

	// -- Check if user is deleting an only admin
	var count int64
	result = handler.DB.Model(&models.User{}).Where("role = ?", "Admin").Count(&count)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}
	if count < 2 && user.Role == "Admin" {
		c.JSON(400, gin.H{
			"error": "cannot delete the only admin",
		})
		return
	}

	// -- Query updated user by ID
	var currentUser models.User
	result = handler.DB.Where("email = ?", email).First(&currentUser)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// -- Prepare for update
	if err := user.PrepareForUpdate(currentUser.ID); err != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// -- Delete user
	user.Deleted = true
	result = handler.DB.Save(&user)
	if result.Error != nil {
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": fmt.Sprintf("user %d deleted", userID),
	})
}
