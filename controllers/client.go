package controllers

import (
	"database/sql"
	"log"
	"server/models"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
)

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
		c.id, c.code, c.name, c.address, c.phone, c.latitude, c.longitude, c.note, smd.id, smd.platform, smd.url 
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
	c.id, c.code, c.name, c.address, c.phone, c.latitude, c.longitude, c.note, smd.id, smd.platform, smd.url 
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

		// -- Apply accounts data without social medias
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

			// -- Append social media to tmpClient
			if tmpClient.Id == scanClient.Id {
				tmpClient.SocialMedias = append(tmpClient.SocialMedias, scanSocialMedia)
				continue
			}

			// -- Append tmpClient to accounts
			if tmpClient.Id != 0 {
				clients = append(clients, tmpClient)
			}

			// -- Set scanClient to tmpClient and append social media
			tmpClient = scanClient
			// -- Social media can be null
			if scanSocialMedia.Id != 0 {
				tmpClient.SocialMedias = append(tmpClient.SocialMedias, scanSocialMedia)
				continue
			}

			// -- Reset social medias to empty array
			tmpClient.SocialMedias = make([]models.SocialMediaData, 0)
		}
	}

	c.JSON(200, utils.NewResponse(200, "success", models.ClientsToResponse(clients)))
}
