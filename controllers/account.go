package controllers

import (
	"errors"
	"log"
	"server/models"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountResponse struct {
	ID             uint   `json:"id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Gender         string `json:"gender"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	Phone          string `json:"phone"`
	SecondaryPhone string `json:"secondary_phone"`
	SocialMedias   []struct {
		ID       uint   `json:"id"`
		Platform string `json:"platform"`
		URL      string `json:"url"`
	} `json:"social_medias"`
}

func (r *AccountResponse) FromModel(account models.Account, socialMedias []models.SocialMedia) {
	r.ID = account.ID
	r.Code = account.Code
	r.Name = account.Name
	r.Email = account.Email
	r.Address = account.Address
	r.Phone = account.Phone
	r.SecondaryPhone = account.SecondaryPhone
	r.SocialMedias = make([]struct {
		ID       uint   `json:"id"`
		Platform string `json:"platform"`
		URL      string `json:"url"`
	}, len(socialMedias))

	// -- Assign social medias
	for i, sm := range socialMedias {
		r.SocialMedias[i].ID = sm.ID
		r.SocialMedias[i].Platform = sm.Platform
		r.SocialMedias[i].URL = sm.URL
	}

	// -- Deserialize
	r.Gender = utils.DeserializeGender(r.Gender)
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

func (handler *AccountHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Query accounts
	var accounts []models.Account
	if err := handler.DB.Limit(paginationQueryParams.Limit).Offset(paginationQueryParams.Offset).Preload("SocialMedias").Find(&accounts).Error; err != nil {
		log.Printf("Error getting accounts from db: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare response
	var response []AccountResponse
	for _, account := range accounts {
		var accountResponse AccountResponse
		accountResponse.FromModel(account, account.SocialMedias)
		response = append(response, accountResponse)
	}

	c.JSON(200, utils.NewResponse(200, "success", response))
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
	result := handler.DB.Preload("SocialMedias").Where("id = ?", account.ID).First(&account)
	if result.Error != nil {
		log.Printf("Error getting account from db: %v\n", result.Error)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	var response AccountResponse
	response.FromModel(account, account.SocialMedias)
	c.JSON(201, utils.NewResponse(201, "account created", response))
}
