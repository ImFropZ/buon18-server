package controllers

import (
	"errors"
	"log"
	"server/models"
	"server/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountResponse struct {
	ID             uint                  `json:"id"`
	Code           string                `json:"code"`
	Name           string                `json:"name"`
	Gender         string                `json:"gender"`
	Email          string                `json:"email"`
	Address        string                `json:"address"`
	Phone          string                `json:"phone"`
	SecondaryPhone string                `json:"secondary_phone"`
	SocialMedias   []SocialMediaResponse `json:"social_medias" gorm:"foreignKey:AccountID"`
}

type SocialMediaResponse struct {
	ID       uint   `json:"id"`
	Platform string `json:"platform"`
	URL      string `json:"url"`

	// -- Foreign key
	AccountID uint `json:"-"`
}

func (SocialMediaResponse) TableName() string {
	return "social_media"
}

type CreateSocialMediaRequest struct {
	Platform string `json:"platform" binding:"required"`
	URL      string `json:"url" binding:"required"`
}

type CreateAccountRequest struct {
	Code           string                     `json:"code" binding:"required"`
	Name           string                     `json:"name" binding:"required"`
	Email          string                     `json:"email"`
	Gender         string                     `json:"gender"`
	Address        string                     `json:"address"`
	Phone          string                     `json:"phone" binding:"required"`
	SecondaryPhone string                     `json:"secondary_phone"`
	SocialMedias   []CreateSocialMediaRequest `json:"social_medias"`
}

type AccountHandler struct {
	DB *gorm.DB
}

func (handler *AccountHandler) First(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user ID. user ID should be an integer"))
		return
	}

	var account AccountResponse
	result := handler.DB.Model(&models.Account{}).Where("id = ?", id).Preload("SocialMedias").First(&account)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, utils.NewErrorResponse(404, "account not found"))
			return
		}

		log.Printf("Error finding account in database: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", account))
}

func (handler *AccountHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Query accounts
	var accounts []AccountResponse
	if err := handler.DB.Model(&models.Account{}).Limit(paginationQueryParams.Limit).Offset(paginationQueryParams.Offset).Preload("SocialMedias").Find(&accounts).Error; err != nil {
		log.Printf("Error getting accounts from db: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", accounts))
}

func (handler *AccountHandler) Create(c *gin.Context) {
	// -- Get email
	email, _ := c.Get("email")
	if email == nil {
		log.Printf("Error getting email from context: %v\n", errors.New("email not found in context"))
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Parse request
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. required fields: code, name, phone"))
		return
	}

	// -- Create account
	account := models.Account{
		Code:  req.Code,
		Name:  req.Name,
		Phone: req.Phone,
	}

	// -- Get current action user
	var user models.User
	if err := handler.DB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("Error getting user from db: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare for create
	if err := account.PrepareForCreate(user.ID, user.ID); err != nil {
		log.Printf("Error preparing account for create: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Add provided fields
	if req.Email != "" {
		account.Email = req.Email
	}
	if req.Address != "" {
		account.Address = req.Address
	}
	if req.SecondaryPhone != "" {
		account.SecondaryPhone = req.SecondaryPhone
	}
	account.Gender = utils.SerializeGender(req.Gender)

	// -- Save account
	if err := handler.DB.Create(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			c.JSON(400, utils.NewErrorResponse(400, "account with either code or phone already exists"))
			return
		}

		log.Printf("Error creating account: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Create social medias
	for _, sm := range req.SocialMedias {
		socialMedia := models.SocialMedia{
			AccountID: account.ID,
			Platform:  strings.ToLower(sm.Platform),
			URL:       sm.URL,
		}
		if err := socialMedia.PrepareForCreate(account.ID, user.ID); err != nil {
			log.Printf("Error preparing social media for create: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
		if err := handler.DB.Create(&socialMedia).Error; err != nil {
			log.Printf("Error creating social media: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Get account from db
	var response AccountResponse
	result := handler.DB.Model(&models.Account{}).Preload("SocialMedias").Where("id = ?", account.ID).First(&response)
	if result.Error != nil {
		log.Printf("Error getting account from db: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(201, utils.NewResponse(201, "account created", response))
}
