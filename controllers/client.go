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
