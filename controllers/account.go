package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"server/database"
	"server/models"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

type AccountResponse struct {
	Id             uint                  `json:"id"`
	Code           string                `json:"code"`
	Name           string                `json:"name"`
	Gender         string                `json:"gender"`
	Email          string                `json:"email"`
	Address        string                `json:"address"`
	Phone          string                `json:"phone"`
	SecondaryPhone string                `json:"secondary_phone"`
	SocialMedias   []SocialMediaResponse `json:"social_medias"`
}

type SocialMediaResponse struct {
	Id       uint   `json:"id"`
	Platform string `json:"platform"`
	URL      string `json:"url"`
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
	DB *sql.DB
}

func (handler *AccountHandler) First(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT 
	a.id, a.code, a.name, a.gender, a.email, a.address, a.phone, a.secondary_phone, COALESCE(sm.id, 0), COALESCE(sm.platform, ''), COALESCE(sm.url, '')
	FROM
		"account" as a
			LEFT JOIN
		"social_media" as sm ON a.id = sm.account_id
	WHERE
		a.id = ?`, id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query database
	var account AccountResponse
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error querying database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		defer rows.Close()

		account.SocialMedias = make([]SocialMediaResponse, 0)
		for rows.Next() {
			var socialMedia SocialMediaResponse
			if err := rows.Scan(&account.Id, &account.Code, &account.Name, &account.Gender, &account.Email, &account.Address, &account.Phone, &account.SecondaryPhone, &socialMedia.Id, &socialMedia.Platform, &socialMedia.URL); err != nil {
				log.Printf("Error scanning row: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}
			if socialMedia.Id != 0 {
				account.SocialMedias = append(account.SocialMedias, socialMedia)
			}
		}
	}

	// -- Return 404 if account.Id == 0 mean account not found
	if account.Id == 0 {
		c.JSON(404, utils.NewErrorResponse(404, "account not found"))
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

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT 
	a.id, a.code, a.name, a.gender, a.email, a.address, a.phone, a.secondary_phone, COALESCE(sm.id, 0), COALESCE(sm.platform, ''), COALESCE(sm.url, '')
	FROM
		"account" as a
			LEFT JOIN
		"social_media" as sm ON a.id = sm.account_id
	ORDER BY a.id 
	LIMIT ? OFFSET ?`, paginationQueryParams.Limit, paginationQueryParams.Offset).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query accounts
	var accounts []AccountResponse
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error querying database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		defer rows.Close()

		// -- Apply accounts data without social medias
		var tmpAccount AccountResponse
		tmpAccount.SocialMedias = make([]SocialMediaResponse, 0)

		for rows.Next() {
			var scanAccount AccountResponse
			var scanSocialMedia SocialMediaResponse
			if err := rows.Scan(&scanAccount.Id, &scanAccount.Code, &scanAccount.Name, &scanAccount.Gender, &scanAccount.Email, &scanAccount.Address, &scanAccount.Phone, &scanAccount.SecondaryPhone, &scanSocialMedia.Id, &scanSocialMedia.Platform, &scanSocialMedia.URL); err != nil {
				log.Printf("Error scanning row: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Append social media to tmpAccount
			if tmpAccount.Id == scanAccount.Id {
				tmpAccount.SocialMedias = append(tmpAccount.SocialMedias, scanSocialMedia)
				continue
			}

			// -- Append tmpAccount to accounts
			if tmpAccount.Id != 0 {
				accounts = append(accounts, tmpAccount)
			}

			// -- Set scanAccount to tmpAccount and append social media
			tmpAccount = scanAccount
			// -- Social media can be null
			if scanSocialMedia.Id != 0 {
				tmpAccount.SocialMedias = append(tmpAccount.SocialMedias, scanSocialMedia)
				continue
			}

			// -- Reset social medias to empty array
			tmpAccount.SocialMedias = make([]SocialMediaResponse, 0)
		}
	}

	c.JSON(200, utils.NewResponse(200, "success", accounts))
}

func (handler *AccountHandler) Create(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, err := c.Get("user_id"); !err {
		log.Printf("Error getting user id: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
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

	// -- Prepare for create
	if err := account.PrepareForCreate(userId, userId); err != nil {
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

	// -- Prepare sql query
	query, params, err := bqb.New(`INSERT INTO 
	"account" (code, name, phone, email, address, secondary_phone, gender, cid, ctime, mid, mtime)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`, account.Code, account.Name, account.Phone, account.Email, account.Address, account.SecondaryPhone, account.Gender, account.CId, account.CTime, account.MId, account.MTime).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Begin transaction
	tx, err := handler.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Create account
	var createdUserId uint
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		if pqErr, ok := row.Err().(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode(database.PQ_ERROR_CODES[database.DUPLICATE]) {
				c.JSON(400, utils.NewErrorResponse(400, "code or phone already exists"))
				return
			}
		}

		log.Printf("Error scaning account: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&createdUserId); err != nil {
			tx.Rollback()
			log.Printf("Error scaning account: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	bqbQuery := bqb.New(`INSERT INTO "social_media" (account_id, platform, url, cid, ctime, mid, mtime) VALUES`)
	socialMedia := models.SocialMedia{}
	if err := socialMedia.PrepareForCreate(userId, userId); err != nil {
		tx.Rollback()
		log.Printf("Error preparing social media for create: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Create social medias
	for _, sm := range req.SocialMedias {
		// -- Append social media to bqb query
		bqbQuery.Space("(?, ?, ?, ?, ?, ?, ?),", createdUserId, sm.Platform, sm.URL, socialMedia.CId, socialMedia.CTime, socialMedia.MId, socialMedia.MTime)
	}

	// -- Prepare social media query
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing social media query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	// -- Remove last comma
	query = query[:len(query)-1]

	// -- Get account from db
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error creating social media: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error commiting transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(201, utils.NewResponse(201, fmt.Sprintf("account %d created", createdUserId), nil))
}
