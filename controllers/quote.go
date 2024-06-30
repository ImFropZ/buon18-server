package controllers

import (
	"database/sql"
	"log"
	"server/models"
	"server/utils"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
)

type QuoteHandler struct {
	DB *sql.DB
}

func (handler *QuoteHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT q.id, q.code, q.date, q.expiry_date, COALESCE(q.note, ''), q.subtotal, q.discount, q.total, q.client_id, q.account_id, q.status, q.cid, qt.id, qt.name, COALESCE(qt.description, ''), qt.quantity, qt.unit_price
	FROM "quote" as q
	LEFT JOIN "quote_item" as qt ON qt.quote_id = q.id 
	ORDER BY q.id, qt.id
	LIMIT ? OFFSET ?`, paginationQueryParams.Limit, paginationQueryParams.Offset).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query quotes from database
	var quotes []models.Quote = make([]models.Quote, 0)
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error finding users in database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		var tmpQuote models.Quote
		tmpQuote.QuoteItems = make([]models.QuoteItem, 0)

		for rows.Next() {
			var scanQuote models.Quote
			var scanQuoteItem models.QuoteItem
			if err := rows.Scan(&scanQuote.Id, &scanQuote.Code, &scanQuote.Date, &scanQuote.ExpiryDate, &scanQuote.Note, &scanQuote.Subtotal, &scanQuote.Discount, &scanQuote.Total, &scanQuote.ClientId, &scanQuote.AccountId, &scanQuote.Status, &scanQuote.CId, &scanQuoteItem.Id, &scanQuoteItem.Name, &scanQuoteItem.Description, &scanQuoteItem.Quantity, &scanQuoteItem.UnitPrice); err != nil {
				log.Printf("Error scanning user from database: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Append social media to tmpQuote's social medias
			if tmpQuote.Id == scanQuote.Id {
				tmpQuote.QuoteItems = append(tmpQuote.QuoteItems, scanQuoteItem)
				continue
			}

			// -- Append tmpQuote to quotes if tmpQuote.Id are not default value
			if tmpQuote.Id != 0 {
				quotes = append(quotes, tmpQuote)
				// -- Reset
				tmpQuote = models.Quote{}
				tmpQuote.QuoteItems = make([]models.QuoteItem, 0)
			}

			// -- Assign scanQuote to tmpQuote
			tmpQuote = scanQuote
			if scanQuoteItem.Id != 0 {
				tmpQuote.QuoteItems = append(tmpQuote.QuoteItems, scanQuoteItem)
			}
		}

		// -- Append last tmpQuote to quotes if tmpQuote.Id are not default value
		if tmpQuote.Id != 0 {
			quotes = append(quotes, tmpQuote)
		}
	}

	c.JSON(200, utils.NewResponse(200, "success", models.QuotesToResponse(quotes)))
}
