package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"server/database"
	"server/models"
	"server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nullism/bqb"
)

type CreateQuoteRequest struct {
	Code       string    `json:"code" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	ExpiryDate time.Time `json:"expiry_date" binding:"required"`
	Note       string    `json:"note"`
	Discount   float64   `json:"discount"`
	ClientId   uint      `json:"client_id" binding:"required"`
	AccountId  uint      `json:"account_id" binding:"required"`
	Items      []struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Quantity    uint    `json:"quantity" binding:"required"`
		UnitPrice   float64 `json:"unit_price" binding:"required"`
	} `json:"items" binding:"required"`
}

type QuoteHandler struct {
	DB *sql.DB
}

func (handler *QuoteHandler) First(c *gin.Context) {
	// -- Get id
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. quote Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT q.id, q.code, q.date, q.expiry_date, COALESCE(q.note, ''), q.subtotal, q.discount, q.total, q.client_id, q.account_id, q.status, q.cid, qt.id, qt.name, COALESCE(qt.description, ''), qt.quantity, qt.unit_price
	FROM "quote" as q
	LEFT JOIN "quote_item" as qt ON qt.quote_id = q.id
	WHERE q.id = ? 
	ORDER BY qt.id`, quoteId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query quote from database
	var quote models.Quote
	quote.QuoteItems = make([]models.QuoteItem, 0)
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error finding quote in database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		var tmpQuoteItem models.QuoteItem
		for rows.Next() {
			var scanQuoteItem models.QuoteItem
			if err := rows.Scan(&quote.Id, &quote.Code, &quote.Date, &quote.ExpiryDate, &quote.Note, &quote.Subtotal, &quote.Discount, &quote.Total, &quote.ClientId, &quote.AccountId, &quote.Status, &quote.CId, &scanQuoteItem.Id, &scanQuoteItem.Name, &scanQuoteItem.Description, &scanQuoteItem.Quantity, &scanQuoteItem.UnitPrice); err != nil {
				log.Printf("Error scanning quote from database: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			// -- Append social media to tmpQuote's social medias
			if tmpQuoteItem.Id == scanQuoteItem.Id {
				tmpQuoteItem = scanQuoteItem
				continue
			}

			// -- Append tmpQuote to quotes if tmpQuote.Id are not default value
			if tmpQuoteItem.Id != 0 {
				quote.QuoteItems = append(quote.QuoteItems, tmpQuoteItem)
			}

			// -- Assign scanQuote to tmpQuote
			tmpQuoteItem = scanQuoteItem
		}

		// -- Append last tmpQuote to quotes if tmpQuote.Id are not default value
		if tmpQuoteItem.Id != 0 {
			quote.QuoteItems = append(quote.QuoteItems, tmpQuoteItem)
		}
	}

	// -- Check if quote is not found
	if quote.Id == 0 {
		c.JSON(404, utils.NewErrorResponse(404, "quote not found"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", quote.ToResponse()))
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

func (handler *QuoteHandler) Create(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Parse request body
	var request CreateQuoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body. required fields: code, date, expiry_date, client_id, account_id, items"))
		return
	}

	// -- Validate items
	if len(request.Items) == 0 {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body. items should not be empty"))
		return
	}

	// -- Prepare quote creation timestamps
	tmpQuote := models.Quote{}
	tmpQuote.PrepareForCreate(userId)

	// -- Calculate subtotal and total
	var subtotal float64
	for _, item := range request.Items {
		subtotal += float64(item.Quantity) * item.UnitPrice
	}
	total := subtotal - request.Discount

	// -- Insert quote
	query, params, err := bqb.New(`INSERT INTO "quote" (code, date, expiry_date, note, discount, subtotal, total, client_id, account_id, status, cid, ctime, mid, mtime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`, request.Code, request.Date, request.ExpiryDate, request.Note, request.Discount, subtotal, total, request.ClientId, request.AccountId, "Draft", tmpQuote.CId, tmpQuote.CTime, tmpQuote.MId, tmpQuote.MTime).ToPgsql()
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

	// -- Query quote from database
	var createdQuoteId uint
	if err := tx.QueryRow(query, params...).Scan(&createdQuoteId); err != nil {
		tx.Rollback()

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode(database.PQ_ERROR_CODES[database.DUPLICATE]) {
				c.JSON(400, utils.NewErrorResponse(400, "quote code already exists"))
				return
			}
		}

		log.Printf("Error inserting quote to database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query
	bqbQuery := bqb.New(`INSERT INTO "quote_item" (quote_id, name, description, quantity, unit_price, cid, ctime, mid, mtime) VALUES `)

	// -- Insert quote items
	for _, item := range request.Items {
		bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?, ?),`, createdQuoteId, item.Name, item.Description, item.Quantity, item.UnitPrice, tmpQuote.CId, tmpQuote.CTime, tmpQuote.MId, tmpQuote.MTime)
	}

	// -- Remove last comma
	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	query = query[:len(query)-1]

	// -- Insert quote items
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error inserting quote items to database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(201, utils.NewResponse(201, fmt.Sprintf("quote %d created successfully", createdQuoteId), nil))
}

func (handler *QuoteHandler) Delete(c *gin.Context) {
	// -- Get id
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. quote Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`DELETE FROM "quote_item" WHERE quote_id = ?`, quoteId).ToPgsql()
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

	// -- Delete quote items from database
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error deleting quote items from database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query
	query, params, err = bqb.New(`DELETE FROM "quote" WHERE id = ?`, quoteId).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Delete quote from database
	if result, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error deleting quote from database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if affected, _ := result.RowsAffected(); affected == 0 {
			tx.Rollback()
			c.JSON(404, utils.NewErrorResponse(404, "quote not found"))
			return
		}
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("quote %d deleted successfully", quoteId), nil))
}

func (handler *QuoteHandler) DeleteItem(c *gin.Context) {
	// -- Get quote id and quote item id
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. quote Id should be an integer"))
		return
	}
	quoteItemId, err := strconv.Atoi(c.Param("qid"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote item Id. quote item Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT COUNT(*) FROM "quote_item" WHERE quote_id = ?`, quoteId).ToPgsql()
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

	// -- Check if quote item are only one left
	var count int
	if err := tx.QueryRow(query, params...).Scan(&count); err != nil {
		tx.Rollback()
		log.Printf("Error counting quote items: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	if count == 1 {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "quote should have at least one item. you can delete the quote instead"))
		return
	}

	// -- Prepare sql query
	query, params, err = bqb.New(`DELETE FROM "quote_item" as qi 
	USING "quote" as q
	WHERE qi.id = ?
	AND q.id = ?`, quoteItemId, quoteId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Delete quote item from database
	if result, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error deleting quote item from database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		if affected, _ := result.RowsAffected(); affected == 0 {
			tx.Rollback()
			c.JSON(404, utils.NewErrorResponse(404, "quote item not found"))
			return
		}
	}

	// -- Prepare sql query (UPDATE QUOTE TOTAL)
	query, params, err = bqb.New(`CALL update_quote_total(?)`, quoteId).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update quote total
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error updating quote total: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("Quote item %d successfully deleted from quote %d", quoteItemId, quoteId), nil))
}
