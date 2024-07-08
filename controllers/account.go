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

type CreateAccountRequest struct {
	Code           string                            `json:"code" binding:"required"`
	Name           string                            `json:"name" binding:"required"`
	Email          string                            `json:"email"`
	Gender         string                            `json:"gender"`
	Address        string                            `json:"address"`
	Phone          string                            `json:"phone" binding:"required"`
	SecondaryPhone string                            `json:"secondary_phone"`
	SocialMedias   []models.CreateSocialMediaRequest `json:"social_medias"`
}

type UpdateAccountRequest struct {
	Code           string                            `json:"code"`
	Name           string                            `json:"name"`
	Email          string                            `json:"email"`
	Gender         string                            `json:"gender"`
	Address        string                            `json:"address"`
	SecondaryPhone string                            `json:"secondary_phone"`
	Phone          string                            `json:"phone"`
	SocialMedias   []models.UpdateSocialMediaRequest `json:"social_medias"`
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
	a.id, a.code, a.name, a.gender, COALESCE(a.email, ''), COALESCE(a.address, ''), a.phone, COALESCE(a.secondary_phone, ''), COALESCE(smd.id, 0), COALESCE(smd.platform, ''), COALESCE(smd.url, '')
	FROM
		"account" as a
			LEFT JOIN
		"social_media_data" as smd ON a.social_media_id = smd.social_media_id
	WHERE
		a.id = ?
	ORDER BY smd.id`, id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query database
	var account models.Account
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error querying database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		defer rows.Close()

		account.SocialMedias = make([]models.SocialMediaData, 0)
		for rows.Next() {
			var socialMedia models.SocialMediaData
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

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"account": account.ToResponse(),
	}))
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
	a.id, a.code, a.name, a.gender, COALESCE(a.email, ''), COALESCE(a.address, ''), a.phone, COALESCE(a.secondary_phone, ''), COALESCE(smd.id, 0), COALESCE(smd.platform, ''), COALESCE(smd.url, '')
	FROM
		"account" as a
			LEFT JOIN
		"social_media_data" as smd ON a.social_media_id = smd.social_media_id
	ORDER BY a.id, smd.id
	LIMIT ? OFFSET ?`, paginationQueryParams.Limit, paginationQueryParams.Offset).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query accounts
	var accounts []models.Account
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error querying database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		defer rows.Close()

		// -- Apply accounts data without social medias
		var tmpAccount models.Account
		tmpAccount.SocialMedias = make([]models.SocialMediaData, 0)

		for rows.Next() {
			var scanAccount models.Account
			var scanSocialMedia models.SocialMediaData
			if err := rows.Scan(&scanAccount.Id, &scanAccount.Code, &scanAccount.Name, &scanAccount.Gender, &scanAccount.Email, &scanAccount.Address, &scanAccount.Phone, &scanAccount.SecondaryPhone, &scanSocialMedia.Id, &scanSocialMedia.Platform, &scanSocialMedia.URL); err != nil {
				log.Printf("Error scanning row: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Append social media to tmpAccount's social medias
			if tmpAccount.Id == scanAccount.Id {
				tmpAccount.SocialMedias = append(tmpAccount.SocialMedias, scanSocialMedia)
				continue
			}

			// -- Append tmpAccount to accounts if tmpAccount.Id are not default value
			if tmpAccount.Id != 0 {
				accounts = append(accounts, tmpAccount)
				// -- Reset
				tmpAccount = models.Account{}
				tmpAccount.SocialMedias = make([]models.SocialMediaData, 0)
			}

			// -- Assign scanAccount to tmpAccount
			tmpAccount = scanAccount
			if scanSocialMedia.Id != 0 {
				tmpAccount.SocialMedias = append(tmpAccount.SocialMedias, scanSocialMedia)
			}
		}

		// -- Append last tmpAccount to accounts
		if tmpAccount.Id != 0 {
			accounts = append(accounts, tmpAccount)
		}
	}

	// -- Count total accounts
	query, params, err = bqb.New(`SELECT COUNT(*) FROM "account"`).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	var total uint
	if err := handler.DB.QueryRow(query, params...).Scan(&total); err != nil {
		log.Printf("Error getting total accounts: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"total":    total,
		"accounts": models.AccountsToResponse(accounts),
	}))
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

	// -- Prepare sql query (CREATE SOCIAL MEDIA)
	query, _, err := bqb.New(`INSERT INTO "social_media" DEFAULT VALUES RETURNING id`).ToPgsql()
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

	// -- Create social media
	var createdSocialMediaId uint
	if row := tx.QueryRow(query); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error creating social media: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&createdSocialMediaId); err != nil {
			tx.Rollback()
			log.Printf("Error scaning social media: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Prepare sql query (CREATE ACCOUNT)
	query, params, err := bqb.New(`INSERT INTO 
	"account" (code, name, phone, email, address, secondary_phone, gender, social_media_id, cid, ctime, mid, mtime)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`, account.Code, account.Name, account.Phone, account.Email, account.Address, account.SecondaryPhone, account.Gender, createdSocialMediaId, account.CId, account.CTime, account.MId, account.MTime).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
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

	// -- Prepare sql query (CREATE SOCIAL MEDIA DATA)
	bqbQuery := bqb.New(`INSERT INTO "social_media_data" (social_media_id, platform, url, cid, ctime, mid, mtime) VALUES`)
	socialMediaData := models.SocialMediaData{}
	if err := socialMediaData.PrepareForCreate(userId, userId); err != nil {
		tx.Rollback()
		log.Printf("Error preparing social media for create: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	for _, sm := range req.SocialMedias {
		// -- Append social media to bqb query
		bqbQuery.Space("(?, ?, ?, ?, ?, ?, ?),", createdSocialMediaId, sm.Platform, sm.URL, socialMediaData.CId, socialMediaData.CTime, socialMediaData.MId, socialMediaData.MTime)
	}

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

func (handler *AccountHandler) Update(c *gin.Context) {
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
	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request"))
		return
	}

	// -- Get account id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query (GET ACCOUNT)
	query, params, err := bqb.New(`SELECT id, code, name, COALESCE(email, ''), gender, COALESCE(address, ''), phone, COALESCE(secondary_phone, ''), social_media_id FROM "account" WHERE id = ?`, id).ToPgsql()
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

	// -- Get account from db
	account := models.Account{}
	var accountSocialMediaId uint
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error getting account: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&account.Id, &account.Code, &account.Name, &account.Email, &account.Gender, &account.Address, &account.Phone, &account.SecondaryPhone, &accountSocialMediaId); err != nil {
			tx.Rollback()
			log.Printf("Error scanning account: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Prepare sql query (UPDATE ACCOUNT)
	bqbQuery := bqb.New(`UPDATE "account" SET`)

	tmpAccount := models.Account{}
	if err := tmpAccount.PrepareForUpdate(userId); err != nil {
		log.Printf("Error preparing account for update: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	if req.Code != "" && req.Code != account.Code {
		bqbQuery.Space(" code = ?,", req.Code)
	}
	if req.Name != "" && req.Name != account.Name {
		bqbQuery.Space(" name = ?,", req.Name)
	}
	if req.Email != "" && req.Email != account.Email {
		bqbQuery.Space(" email = ?,", req.Email)
	}
	if req.Gender != "" {
		if g := utils.SerializeGender(req.Gender); g != account.Gender {
			bqbQuery.Space(` gender = ?,`, g)
		}
	}
	if req.Address != "" && req.Address != account.Address {
		bqbQuery.Space(" address = ?,", req.Address)
	}
	if req.Phone != "" && req.Phone != account.Phone {
		bqbQuery.Space(" phone = ?,", req.Phone)
	}
	if req.SecondaryPhone != "" && req.SecondaryPhone != account.SecondaryPhone {
		bqbQuery.Space(" secondary_phone = ?,", req.SecondaryPhone)
	}

	// -- Append mid and mtime
	query, params, err = bqbQuery.Space(" mid = ?, mtime = ? WHERE id = ? RETURNING id", tmpAccount.MId, tmpAccount.MTime, id).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update account
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error updating account: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update social medias
	if len(req.SocialMedias) == 0 {
		// -- Commit transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error commiting transaction: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		c.JSON(200, utils.NewResponse(200, fmt.Sprintf("account %d updated", id), nil))
		return
	}

	// -- Prepare for update social media and can be used for create social media
	tmpSocialMedia := models.SocialMediaData{}
	if err := tmpSocialMedia.PrepareForCreate(userId, userId); err != nil {
		tx.Rollback()
		log.Printf("Error preparing social media for update: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Separate social medias to create and update
	var createSocialMedias []models.SocialMediaData
	var updateSocialMedias []models.SocialMediaData
	for _, sm := range req.SocialMedias {
		if sm.Id == 0 {
			createSocialMedias = append(createSocialMedias, models.SocialMediaData{
				Platform: sm.Platform,
				URL:      sm.URL,
			})
			continue
		}
		updateSocialMedias = append(updateSocialMedias, models.SocialMediaData{
			Id:       sm.Id,
			Platform: sm.Platform,
			URL:      sm.URL,
		})
	}

	// -- Update social medias
	for _, sm := range updateSocialMedias {
		if sm.Platform == "" && sm.URL == "" {
			continue
		}
		bqbQuery = bqb.New(`UPDATE "social_media_data" SET`)
		if sm.Platform != "" {
			bqbQuery.Space("platform = ?,", sm.Platform)
		}
		if sm.URL != "" {
			bqbQuery.Space("url = ?,", sm.URL)
		}
		query, params, err = bqbQuery.Space("mid = ?, mtime = ? WHERE id = ? AND social_media_id = ?", tmpSocialMedia.MId, tmpSocialMedia.MTime, sm.Id, accountSocialMediaId).ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing sql query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Update social media
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error updating social media: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	bqbQuery = bqb.New(`INSERT INTO "social_media_data" (social_media_id, platform, url, cid, ctime, mid, mtime) VALUES`)

	// -- Create social medias
	var validCreateCount uint
	for _, sm := range createSocialMedias {
		// -- Need both fields to create social media
		if sm.Platform == "" || sm.URL == "" {
			tx.Rollback()
			c.JSON(400, utils.NewErrorResponse(400, "invalid social media. required fields: platform, url"))
			return
		}
		bqbQuery.Space("(?, ?, ?, ?, ?, ?, ?),", accountSocialMediaId, sm.Platform, sm.URL, tmpSocialMedia.CId, tmpSocialMedia.CTime, tmpSocialMedia.MId, tmpSocialMedia.MTime)
		validCreateCount++
	}

	// -- Prepare social media query
	if validCreateCount > 0 {
		query, params, err = bqbQuery.ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing social media query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Remove last comma
		query = query[:len(query)-1]

		// -- Create social medias
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error creating social media: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error commiting transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("account %d updated", id), nil))
}

func (handler *AccountHandler) DeleteSocialMedia(c *gin.Context) {
	// -- Get id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Get social media id
	socialMediaId, err := strconv.Atoi(c.Param("smid"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid social media Id. social media Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`DELETE FROM "social_media_data" as smd
	USING "account" as a
	WHERE smd.social_media_id = a.social_media_id
	AND smd.id = ?
	AND a.id = ?`, socialMediaId, id).ToPgsql()
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

	// -- Delete social media
	if result, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error deleting social media: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if n, err := result.RowsAffected(); err != nil {
			tx.Rollback()
			log.Printf("Error getting rows affected: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		} else if n == 0 {
			tx.Rollback()
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("social media %d not found from account %d", socialMediaId, id)))
			return
		}
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error commiting transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("social media %d deleted", socialMediaId), nil))
}

func (handler *AccountHandler) Delete(c *gin.Context) {
	// -- Get id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query to delete social medias
	query, params, err := bqb.New(`DELETE FROM "social_media_data" as smd
	USING "account" AS a 
	WHERE smd.social_media_id = a.social_media_id
	AND a.id = ?`, id).ToPgsql()
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

	// -- Delete social medias
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error deleting social medias: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query to delete account
	query, params, err = bqb.New(`DELETE FROM "account" 
	WHERE id = ?
	RETURNING social_media_id`, id).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Delete account
	var targetSocialMediaId uint
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error deleting account: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&targetSocialMediaId); err != nil {
			tx.Rollback()
			log.Printf("Error scaning account: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Prepare sql query to delete social media
	query, params, err = bqb.New(`DELETE FROM "social_media" WHERE id = ?`, targetSocialMediaId).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Delete social media
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error deleting social media: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error commiting transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("account %d deleted", id), nil))
}
