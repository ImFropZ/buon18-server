package controllers

import (
	"errors"
	"fmt"
	"log"
	"server/models"
	"server/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Deleted  string `json:"deleted"`
}

type UserHandler struct {
	DB *gorm.DB
}

func (handler *UserHandler) First(c *gin.Context) {
	var user UserResponse
	result := handler.DB.Model(&models.User{}).First(&user)
	if result.Error != nil {
		log.Printf("Error finding users in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", user))
}

func (handler *UserHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	var users []UserResponse
	result := handler.DB.Model(&models.User{}).Limit(paginationQueryParams.Limit).Offset(paginationQueryParams.Offset).Find(&users)
	if result.Error != nil {
		log.Printf("Error finding users in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", users))
}

func (handler *UserHandler) Create(c *gin.Context) {
	// -- Get email
	email, _ := c.Get("email")
	if email == nil {
		log.Printf("Error getting email from context: %v\n", errors.New("email not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Bind request
	var request CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain name, email, password, and role fields"))
		return
	}

	// -- Validate role
	role, ok := utils.ValidateRole(request.Role)
	if !ok {
		c.JSON(400, utils.NewErrorResponse(400, "invalid role"))
		return
	}

	// -- Hash password
	hashedPwd, err := utils.HashPwd(request.Password)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
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
		log.Printf("Error finding matched email in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare for create
	if err := user.PrepareForCreate(existingUser.ID, existingUser.ID); err != nil {
		log.Printf("Error preparing create fields for user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Create user
	result = handler.DB.Create(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(400, utils.NewErrorResponse(400, "email already exists"))
			return
		}
		log.Printf("Error creating new user: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(201, utils.NewResponse(201, "user created", &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}))
}

func (handler *UserHandler) Update(c *gin.Context) {
	// -- Get email
	email, _ := c.Get("email")
	if email == nil {
		log.Printf("Error getting email from context: %v\n", errors.New("email not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get user ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user ID. user ID should be an integer"))
		return
	}

	// -- Bind request
	var request UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. request should contain either one of name, email, password, and role fields"))
		return
	}

	// -- Validate role
	var role string
	if request.Role != "" {
		if roleStr, ok := utils.ValidateRole(request.Role); !ok {
			c.JSON(400, utils.NewErrorResponse(400, "invalid role"))
			return
		} else {
			role = roleStr
		}
	}

	// -- Query user by ID
	var user models.User
	result := handler.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		c.JSON(404, utils.NewErrorResponse(404, "user not found"))
		return
	}

	// -- Check if user already deleted
	if user.Deleted {
		if lower := strings.ToLower(request.Deleted); lower != "false" || lower == "true" {
			c.JSON(400, utils.NewErrorResponse(400, "user already deleted"))
			return
		}
	}

	// -- Query updated user by email
	var currentUser models.User
	result = handler.DB.Where("email = ?", email).First(&currentUser)
	if result.Error != nil {
		log.Printf("Error finding matched email in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare for update
	if err := user.PrepareForUpdate(currentUser.ID); err != nil {
		log.Printf("Error preparing update fields for user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if user is updating an only admin
	var count int64
	result = handler.DB.Model(&models.User{}).Where("role = ?", "Admin").Count(&count)
	if result.Error != nil {
		log.Printf("Error finding admin role in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	if count < 2 && user.Role == "Admin" && role != "Admin" {
		c.JSON(400, utils.NewErrorResponse(400, "cannot update the only admin"))
		return
	}

	// -- Update user
	if request.Name != "" {
		user.Name = request.Name
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.Password != "" {
		hashedPwd, err := utils.HashPwd(request.Password)
		if err != nil {
			log.Printf("Error hashing password: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
		user.Pwd = hashedPwd
	}
	if role != "" {
		user.Role = role
	}
	if lower := strings.ToLower(request.Deleted); lower == "true" || lower == "false" {
		user.Deleted = lower == "true"
	}
	result = handler.DB.Save(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(400, utils.NewErrorResponse(400, "email already exists"))
			return
		}
		log.Printf("Error saving into database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("user %d updated", user.ID), &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}))
}

func (handler *UserHandler) Delete(c *gin.Context) {
	// -- Get email
	email, _ := c.Get("email")
	if email == nil {
		log.Printf("Error getting email from context: %v\n", errors.New("email not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get user ID
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user ID. user ID should be an integer"))
		return
	}

	// -- Query user by ID
	var user models.User
	result := handler.DB.Where("id = ?", userID).First(&user)
	if result.Error != nil {
		c.JSON(404, utils.NewErrorResponse(404, "user not found"))
		return
	}

	// -- Check if user is deleting itself
	if user.Email == email {
		c.JSON(400, utils.NewErrorResponse(400, "user cannot delete itself"))
		return
	}

	// -- Check if user already deleted
	if user.Deleted {
		c.JSON(400, utils.NewErrorResponse(400, "user already deleted"))
		return
	}

	// -- Check if user is deleting an only admin
	var count int64
	result = handler.DB.Model(&models.User{}).Where("role = ?", "Admin").Count(&count)
	if result.Error != nil {
		log.Printf("Error finding admin role in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	if count < 2 && user.Role == "Admin" {
		c.JSON(400, utils.NewErrorResponse(400, "cannot delete the only admin"))
		return
	}

	// -- Query updated user by ID
	var currentUser models.User
	result = handler.DB.Where("email = ?", email).First(&currentUser)
	if result.Error != nil {
		log.Printf("Error finding matched email in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare for update
	if err := user.PrepareForUpdate(currentUser.ID); err != nil {
		log.Printf("Error preparing update fields for user: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Delete user
	user.Deleted = true
	result = handler.DB.Save(&user)
	if result.Error != nil {
		log.Printf("Error saving into database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("user %d deleted", user.ID), nil))
}
