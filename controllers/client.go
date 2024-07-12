package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"server/database"
	"server/models"
	"server/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

func prepareClientQuery(c *gin.Context, bqbQuery *bqb.Query) {
	// -- Apply query params
	bqbQuery.Space("WHERE")
	if str, ok := c.GetQuery("name_ilike"); ok {
		bqbQuery.Space(`c.name ILIKE ? AND`, "%"+str+"%")
	}
	if str, ok := c.GetQuery("phone_ilike"); ok {
		bqbQuery.Space(`c.phone ILIKE ? AND`, "%"+str+"%")
	}
	if str, ok := c.GetQuery("latitude_range"); ok {
		if value := strings.Split(str, ","); len(value) == 2 {
			if min, err := strconv.ParseFloat(value[0], 64); err == nil {
				if max, err := strconv.ParseFloat(value[1], 64); err == nil {
					bqbQuery.Space(`c.latitude BETWEEN ? AND ? AND`, min, max)
				}
			}
		}
	}
	if str, ok := c.GetQuery("longitude_range"); ok {
		if value := strings.Split(str, ","); len(value) == 2 {
			if min, err := strconv.ParseFloat(value[0], 64); err == nil {
				if max, err := strconv.ParseFloat(value[1], 64); err == nil {
					bqbQuery.Space(`c.longitude BETWEEN ? AND ? AND`, min, max)
				}
			}
		}
	}

	// -- Remove last AND or WHERE
	if strings.HasSuffix(bqbQuery.Parts[len(bqbQuery.Parts)-1].Text, "WHERE") {
		bqbQuery.Parts = bqbQuery.Parts[:len(bqbQuery.Parts)-1]
	} else if strings.HasSuffix(bqbQuery.Parts[len(bqbQuery.Parts)-1].Text, "AND") {
		text := bqbQuery.Parts[len(bqbQuery.Parts)-1].Text
		arr := strings.Split(text, " ")

		bqbQuery.Parts[len(bqbQuery.Parts)-1].Text = strings.Join(arr[:len(arr)-1], " ")
	}
}

type CreateClientRequest struct {
	Code         string                            `json:"code" binding:"required"`
	Name         string                            `json:"name" binding:"required"`
	Address      string                            `json:"address"`
	Phone        string                            `json:"phone" binding:"required"`
	Latitude     float64                           `json:"latitude" binding:"required"`
	Longitude    float64                           `json:"longitude" binding:"required"`
	Note         string                            `json:"note"`
	SocialMedias []models.CreateSocialMediaRequest `json:"social_medias"`
}

type UpdateClientRequest struct {
	Code         string                            `json:"code"`
	Name         string                            `json:"name"`
	Address      string                            `json:"address"`
	Phone        string                            `json:"phone"`
	Latitude     float64                           `json:"latitude"`
	Longitude    float64                           `json:"longitude"`
	Note         string                            `json:"note"`
	SocialMedias []models.UpdateSocialMediaRequest `json:"social_medias"`

	// -- Delete social medias
	DeleteSocialMedias []uint `json:"delete_social_media_ids"`
}

type ClientHandler struct {
	DB *sql.DB
}

func (handler *ClientHandler) First(c *gin.Context) {
	// -- Get id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query (GET CLIENT)
	query, params, err := bqb.New(`SELECT 
		c.id, c.code, c.name, COALESCE(c.address, ''), c.phone, c.latitude, c.longitude, COALESCE(c.note, ''), COALESCE(smd.id, 0), COALESCE(smd.platform, ''), COALESCE(smd.url, '') 
		FROM 
			"client" as c
				LEFT JOIN 
			"social_media_data" as smd ON c.social_media_id = smd.social_media_id
		WHERE c.id = ?
		ORDER BY smd.id`, id).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get client from db
	var client models.Client
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error querying database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		defer rows.Close()

		if rows.Next() {
			var socialMedia models.SocialMediaData
			if err := rows.Scan(&client.Id, &client.Code, &client.Name, &client.Address, &client.Phone, &client.Latitude, &client.Longitude, &client.Note, &socialMedia.Id, &socialMedia.Platform, &socialMedia.URL); err != nil {
				log.Printf("Error scanning row: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}
			if socialMedia.Id != 0 {
				client.SocialMedias = append(client.SocialMedias, socialMedia)
			}
		}
	}

	// -- Check if client exists
	if client.Id == 0 {
		c.JSON(404, utils.NewErrorResponse(404, "client not found"))
		return
	}

	// -- Return client
	c.JSON(200, utils.NewResponse(200, "", gin.H{
		"client": client.ToResponse(),
	}))
}

func (handler *ClientHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Prepare sql query (GET CLIENTS)
	bqbQuery := bqb.New(`SELECT 
	c.id, c.code, c.name, COALESCE(c.address, ''), c.phone, c.latitude, c.longitude, COALESCE(c.note, ''), COALESCE(smd.id, 0), COALESCE(smd.platform, ''), COALESCE(smd.url, '') 
	FROM 
		"client" as c
			LEFT JOIN 
		"social_media_data" as smd ON c.social_media_id = smd.social_media_id`)

	// -- Apply query params
	prepareClientQuery(c, bqbQuery)

	// -- Complete query
	bqbQuery.Space("ORDER BY c.id, smd.id OFFSET ? LIMIT ?", paginationQueryParams.Offset, paginationQueryParams.Limit)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get clients from db
	clients := make([]models.Client, 0)
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error querying database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		defer rows.Close()

		// -- Apply clients data without social medias
		var tmpClient models.Client
		tmpClient.SocialMedias = make([]models.SocialMediaData, 0)

		for rows.Next() {
			var scanClient models.Client
			var scanSocialMedia models.SocialMediaData
			if err := rows.Scan(&scanClient.Id, &scanClient.Code, &scanClient.Name, &scanClient.Address, &scanClient.Phone, &scanClient.Latitude, &scanClient.Longitude, &scanClient.Note, &scanSocialMedia.Id, &scanSocialMedia.Platform, &scanSocialMedia.URL); err != nil {
				log.Printf("Error scanning row: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Append social media to tmpClient's social medias
			if tmpClient.Id == scanClient.Id {
				tmpClient.SocialMedias = append(tmpClient.SocialMedias, scanSocialMedia)
				continue
			}

			// -- Append tmpClient to clients if tmpClient.Id are not default value
			if tmpClient.Id != 0 {
				clients = append(clients, tmpClient)
				// -- Reset
				tmpClient = models.Client{}
				tmpClient.SocialMedias = make([]models.SocialMediaData, 0)
			}

			// -- Assign scanClient to tmpClient
			tmpClient = scanClient
			if scanSocialMedia.Id != 0 {
				tmpClient.SocialMedias = append(tmpClient.SocialMedias, scanSocialMedia)
			}
		}

		// -- Append last tmpClient to clients
		if tmpClient.Id != 0 {
			clients = append(clients, tmpClient)
		}
	}

	// -- Count total clients
	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "client" as c`)

	prepareClientQuery(c, bqbQuery)

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	var total uint
	if err := handler.DB.QueryRow(query, params...).Scan(&total); err != nil {
		log.Printf("Error getting total clients: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"total":   total,
		"clients": models.ClientsToResponse(clients),
	}))
}

func (handler *ClientHandler) Create(c *gin.Context) {
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
	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error binding request: %v\n", err)
		c.JSON(400, utils.NewErrorResponse(400, "invalid request. required fields: code, name, phone, latitude, and longitude. The latitude and longitude different from 0.0"))
		return
	}

	// -- Create client
	client := models.Client{
		Code:      req.Code,
		Name:      req.Name,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// -- Prepare for create
	if err := client.PrepareForCreate(userId, userId); err != nil {
		log.Printf("Error preparing client for create: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Add provided fields
	if req.Address != "" {
		client.Address = req.Address
	}
	if req.Note != "" {
		client.Note = req.Note
	}

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

	// -- Prepare sql query (CREATE CLIENT)
	query, params, err := bqb.New(`INSERT INTO 
	"client" (code, name, phone, latitude, longitude, note, social_media_id, cid, ctime, mid, mtime)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING id`, client.Code, client.Name, client.Phone, client.Latitude, client.Longitude, client.Note, createdSocialMediaId, client.CId, client.CTime, client.MId, client.MTime).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Create client
	var createdClientId uint
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		if pqErr, ok := row.Err().(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode(database.PQ_ERROR_CODES[database.DUPLICATE]) {
				c.JSON(400, utils.NewErrorResponse(400, "code or phone already exists"))
				return
			}
		}

		log.Printf("Error create client: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&createdClientId); err != nil {
			tx.Rollback()
			log.Printf("Error scaning client: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	if len(req.SocialMedias) != 0 {
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

		// -- Get client from db
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

	c.JSON(201, utils.NewResponse(201, fmt.Sprintf("client %d created", createdClientId), nil))
}

func (handler *ClientHandler) Update(c *gin.Context) {
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
	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request"))
		return
	}

	// -- Get id
	targetClientId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query (GET CLIENT)
	query, params, err := bqb.New(`SELECT id, code, name, COALESCE(address, ''), phone, latitude, longitude, COALESCE(note, ''), social_media_id FROM "client" WHERE id = ?`, targetClientId).ToPgsql()
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

	// -- Get client from db
	client := models.Client{}
	var clientSocialMediaId uint
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error getting client: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&client.Id, &client.Code, &client.Name, &client.Address, &client.Phone, &client.Latitude, &client.Longitude, &client.Note, &clientSocialMediaId); err != nil {
			tx.Rollback()

			if err == sql.ErrNoRows {
				c.JSON(404, utils.NewErrorResponse(404, "client not found"))
				return
			}

			log.Printf("Error scanning client: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Prepare sql query (UPDATE CLIENT)
	bqbQuery := bqb.New(`UPDATE "client" SET`)

	tmpClient := models.Client{}
	if err := tmpClient.PrepareForUpdate(userId); err != nil {
		log.Printf("Error preparing client for update: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	if req.Code != "" && req.Code != client.Code {
		bqbQuery.Space(" code = ?,", req.Code)
	}
	if req.Name != "" && req.Name != client.Name {
		bqbQuery.Space(" name = ?,", req.Name)
	}
	if req.Address != "" && req.Address != client.Address {
		bqbQuery.Space(" address = ?,", req.Address)
	}
	if req.Phone != "" && req.Phone != client.Phone {
		bqbQuery.Space(" phone = ?,", req.Phone)
	}
	if req.Latitude != 0.0 && req.Latitude != client.Latitude {
		bqbQuery.Space(" latitude = ?,", req.Latitude)
	}
	if req.Longitude != 0.0 && req.Longitude != client.Longitude {
		bqbQuery.Space(" longitude = ?,", req.Longitude)
	}
	if req.Note != "" && req.Note != client.Note {
		bqbQuery.Space(" note = ?,", req.Note)
	}

	// -- Append mid and mtime
	query, params, err = bqbQuery.Space(" mid = ?, mtime = ? WHERE id = ? RETURNING id", tmpClient.MId, tmpClient.MTime, targetClientId).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update client
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error updating client: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update social medias
	if len(req.SocialMedias) == 0 && len(req.DeleteSocialMedias) == 0 {
		// -- Commit transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Error commiting transaction: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		c.JSON(200, utils.NewResponse(200, fmt.Sprintf("client %d updated", targetClientId), nil))
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
		query, params, err = bqbQuery.Space("mid = ?, mtime = ? WHERE id = ? AND social_media_id = ?", tmpSocialMedia.MId, tmpSocialMedia.MTime, sm.Id, clientSocialMediaId).ToPgsql()
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
		bqbQuery.Space("(?, ?, ?, ?, ?, ?, ?),", clientSocialMediaId, sm.Platform, sm.URL, tmpSocialMedia.CId, tmpSocialMedia.CTime, tmpSocialMedia.MId, tmpSocialMedia.MTime)
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

	// -- Delete social medias
	if len(req.DeleteSocialMedias) > 0 {
		// -- Prepare sql query
		bqbQuery = bqb.New(`DELETE FROM "social_media_data" WHERE id IN (`)
		for _, smid := range req.DeleteSocialMedias {
			bqbQuery.Space("?,", smid)
		}

		query, params, err = bqbQuery.ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing social media query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Remove last comma and add closing bracket
		query = query[:len(query)-1] + ")"

		// -- Delete social medias
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error deleting social media: %v\n", err)
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

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("client %d updated", targetClientId), nil))
}

func (handler *ClientHandler) Delete(c *gin.Context) {
	// -- Get id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid user Id. user Id should be an integer"))
		return
	}

	// -- Prepare sql query to delete social medias
	query, params, err := bqb.New(`DELETE FROM "social_media_data" as smd
	USING "client" AS c 
	WHERE smd.social_media_id = c.social_media_id
	AND c.id = ?`, id).ToPgsql()
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

	// -- Prepare sql query to delete client
	query, params, err = bqb.New(`DELETE FROM "client" 
	WHERE id = ?
	RETURNING social_media_id`, id).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Delete client
	var targetSocialMediaId uint
	if row := tx.QueryRow(query, params...); row.Err() != nil {
		tx.Rollback()
		log.Printf("Error deleting client: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&targetSocialMediaId); err != nil {
			tx.Rollback()

			if err == sql.ErrNoRows {
				c.JSON(404, utils.NewErrorResponse(404, "client not found"))
				return
			}

			log.Printf("Error scaning client: %v\n", err)
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

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("client %d deleted", id), nil))
}
