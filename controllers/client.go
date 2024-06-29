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
	c.JSON(200, utils.NewResponse(200, "", client.ToResponse()))
}

func (handler *ClientHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Prepare sql query (GET CLIENTS)
	query, params, err := bqb.New(`SELECT 
	c.id, c.code, c.name, COALESCE(c.address, ''), c.phone, c.latitude, c.longitude, COALESCE(c.note, ''), COALESCE(smd.id, 0), COALESCE(smd.platform, ''), COALESCE(smd.url, '') 
	FROM 
		"client" as c
			LEFT JOIN 
		"social_media_data" as smd ON c.social_media_id = smd.social_media_id
	ORDER BY c.id, smd.id
	LIMIT ? OFFSET ?`, paginationQueryParams.Limit, paginationQueryParams.Offset).ToPgsql()
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

	c.JSON(200, utils.NewResponse(200, "success", models.ClientsToResponse(clients)))
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

	// -- Create account
	client := models.Client{
		Code:      req.Code,
		Name:      req.Name,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// -- Prepare for create
	if err := client.PrepareForCreate(userId, userId); err != nil {
		log.Printf("Error preparing account for create: %v\n", err)
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

		log.Printf("Error create account: %v\n", row.Err())
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if err := row.Scan(&createdClientId); err != nil {
			tx.Rollback()
			log.Printf("Error scaning account: %v\n", err)
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

		// -- Get account from db
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
