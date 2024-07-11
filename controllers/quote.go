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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/nullism/bqb"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

type UpdateQuoteRequest struct {
	Code       string    `json:"code"`
	Date       time.Time `json:"date"`
	ExpiryDate time.Time `json:"expiry_date"`
	Note       string    `json:"note"`     // -- WARN: Note = "no" means no note
	Discount   float64   `json:"discount"` // -- WARN: Discount = -1 means 0 discount since 0 value means default value
	ClientId   uint      `json:"client_id"`
	AccountId  uint      `json:"account_id"`
	Items      []struct {
		Id          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Quantity    uint    `json:"quantity"`
		UnitPrice   float64 `json:"unit_price"`
	} `json:"items"`

	// -- Delete items
	DeleteItems []uint `json:"delete_item_ids"`
}

type CreateSalesOrderRequest struct {
	Code         string    `json:"code" binding:"required"`
	DeliveryDate time.Time `json:"delivery_date" binding:"required"`
	Note         string    `json:"note"`
}

type UpdateQuoteStatusRequest struct {
	Action string `json:"action" binding:"required"`
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
	query, params, err := bqb.New(`SELECT q.id, q.code, q.date, q.expiry_date, COALESCE(q.note, ''), q.subtotal, q.discount, q.total, q.client_id, q.account_id, q.status, q.cid, COALESCE(qt.id, 0), COALESCE(qt.name, ''), COALESCE(qt.description, ''), COALESCE(qt.quantity, 0), COALESCE(qt.unit_price, 0)
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

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"quote": quote.ToResponse(),
	}))
}

func (handler *QuoteHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Prepare sql query
	bqbQuery := bqb.New(`SELECT q.id, q.code, q.date, q.expiry_date, COALESCE(q.note, ''), q.subtotal, q.discount, q.total, q.client_id, q.account_id, q.status, q.cid, COALESCE(qt.id, 0), COALESCE(qt.name, ''), COALESCE(qt.description, ''), COALESCE(qt.quantity, 0), COALESCE(qt.unit_price, 0)
	FROM "quote" as q
	LEFT JOIN "quote_item" as qt ON qt.quote_id = q.id`)

	// -- Add query if exists
	if paginationQueryParams.Query != "" {
		bqbQuery.Space(`WHERE q.code ILIKE ?`, "%"+paginationQueryParams.Query+"%")
	}

	// -- Complete query
	bqbQuery.Space("ORDER BY q.id, qt.id OFFSET ? LIMIT ?", paginationQueryParams.Offset, paginationQueryParams.Limit)

	query, params, err := bqbQuery.ToPgsql()
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

	// -- Count total quotes
	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "quote"`)

	if paginationQueryParams.Query != "" {
		bqbQuery.Space(`WHERE code ILIKE ?`, "%"+paginationQueryParams.Query+"%")
	}

	query, params, err = bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	var total uint
	if err := handler.DB.QueryRow(query, params...).Scan(&total); err != nil {
		log.Printf("Error getting total quotes: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"total":  total,
		"quotes": models.QuotesToResponse(quotes),
	}))
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

	// -- Validate date
	if !request.ExpiryDate.After(request.Date) {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body. expiry_date should be after date"))
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

			if pqErr.Code == pq.ErrorCode(database.PQ_ERROR_CODES[database.VIOLATE_FOREIGN_KEY]) {
				c.JSON(400, utils.NewErrorResponse(400, "client or account not found"))
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

func (handler *QuoteHandler) UpdateStatus(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Get id
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. quote Id should be an integer"))
		return
	}

	// -- Parse request body
	var request UpdateQuoteStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body. required fields: action"))
		return
	}

	// -- Validate action
	caser := cases.Title(language.English)
	if action := caser.String(request.Action); action == "Sent" || action == "Accept" || action == "Reject" {
		request.Action = action
	} else {
		c.JSON(400, utils.NewErrorResponse(400, "invalid action. action should be sent, accept, or reject"))
		return
	}

	// -- Prepare quote update timestamps
	tmpQuote := models.Quote{}
	tmpQuote.PrepareForUpdate(userId)

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT expiry_date, status FROM "quote" WHERE id = ?`, quoteId).ToPgsql()
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

	// -- Check if quote exists
	var targetQuote models.Quote
	if err := tx.QueryRow(query, params...).Scan(&targetQuote.ExpiryDate, &targetQuote.Status); err != nil {
		tx.Rollback()

		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, "quote not found"))
			return
		}

		log.Printf("Error counting quote: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if status is Expired
	if isExpired(targetQuote.ExpiryDate) {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "quote status is already expired. create a new quote instead"))
		return
	}

	// -- Check if status is already in the sales order
	if targetQuote.Status == "Accept" {
		query, params, err = bqb.New(`SELECT COUNT(*) FROM "sales_order" WHERE quote_id = ?`, quoteId).ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing sql query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Check if quote is already in the sales order
		var count int
		if err := tx.QueryRow(query, params...).Scan(&count); err != nil {
			tx.Rollback()
			log.Printf("Error counting sales order: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		if count > 0 {
			tx.Rollback()
			c.JSON(400, utils.NewErrorResponse(400, "quote status is already in the sales order. you can update status in the sales order instead"))
			return
		}
	}

	// -- Check if request status is Sent and status is Accept or Reject
	if (targetQuote.Status == "Accept" || targetQuote.Status == "Reject") && request.Action == "Sent" {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "quote status is already "+strings.ToLower(targetQuote.Status)+". you can only update to either approved or rejected status"))
		return
	}

	// -- Check if status is Draft (Draft to Sent)
	if targetQuote.Status == "Draft" && request.Action == "Sent" {
		// -- Update status to sent
		query, params, err = bqb.New(`UPDATE "quote" SET status = ?, mid = ?, mtime = ? WHERE id = ?`, request.Action, tmpQuote.MId, tmpQuote.MTime, quoteId).ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing sql query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Update quote
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error updating quote: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		tx.Commit()
		c.JSON(200, utils.NewResponse(200, fmt.Sprintf("quote %d updated successfully", quoteId), nil))
		return
	}

	// -- Prepare sql query
	query, params, err = bqb.New(`UPDATE "quote" SET status = ?, mid = ?, mtime = ? WHERE id = ?`, request.Action, tmpQuote.MId, tmpQuote.MTime, quoteId).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update quote
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()
		log.Printf("Error updating quote: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	tx.Commit()
	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("quote %d updated successfully", quoteId), nil))
}

func (handler *QuoteHandler) CreateSalesOrder(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Get id
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. quote Id should be an integer"))
		return
	}

	// -- Parse request body
	var request CreateSalesOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body. required fields: code, delivery_date"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT status FROM "quote" WHERE id = ?`, quoteId).ToPgsql()
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

	// -- Check if quote exists
	var quoteStatus string
	if err := tx.QueryRow(query, params...).Scan(&quoteStatus); err != nil {
		tx.Rollback()

		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, "quote not found"))
			return
		}

		log.Printf("Error counting quote: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if status is Accept
	if quoteStatus != "Accept" {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "quote status is not accepted"))
		return
	}

	// -- Prepare sales order creation timestamps
	tmpSalesOrder := models.SalesOrder{}
	tmpSalesOrder.PrepareForCreate(userId)

	// -- Prepare sql query
	query, params, err = bqb.New(`INSERT INTO "sales_order" 
	(code, accept_date, delivery_date, note, quote_id, cid, ctime, mid, mtime) 
	VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?)`, request.Code, time.Now(), request.DeliveryDate, request.Note, quoteId, tmpSalesOrder.CId, tmpSalesOrder.CTime, tmpSalesOrder.MId, tmpSalesOrder.MTime).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	// -- Insert sales order
	if _, err := tx.Exec(query, params...); err != nil {
		tx.Rollback()

		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == pq.ErrorCode(database.PQ_ERROR_CODES[database.DUPLICATE]) {
				c.JSON(400, utils.NewErrorResponse(400, "sales order code already exists"))
				return
			}
		}

		log.Printf("Error inserting sales order to database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(201, utils.NewResponse(201, "sales order created successfully", nil))
}

func (handler *QuoteHandler) Update(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Get id
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. quote Id should be an integer"))
		return
	}

	// -- Parse request body
	var request UpdateQuoteRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body."))
		return
	}

	// -- Prepare quote update timestamps
	tmpQuote := models.Quote{}
	tmpQuote.PrepareForUpdate(userId)

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT code, date, expiry_date, COALESCE(note, ''), discount, client_id, account_id, status FROM "quote" WHERE id = ?`, quoteId).ToPgsql()
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

	// -- Check if quote exists
	var quote models.Quote
	if err := tx.QueryRow(query, params...).Scan(&quote.Code, &quote.Date, &quote.ExpiryDate, &quote.Note, &quote.Discount, &quote.ClientId, &quote.AccountId, &quote.Status); err != nil {
		tx.Rollback()

		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, "quote not found"))
			return
		}

		log.Printf("Error counting quote: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Check if status is Accept, Reject, or Expired
	if quote.Status == "Accept" || quote.Status == "Reject" || isExpired(quote.ExpiryDate) {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "quote status is already "+strings.ToLower(quote.Status)))
		return
	}

	// -- Prepare sql query
	bqbQuery := bqb.New(`UPDATE "quote" SET`)

	// -- Update code
	var isDiscountUpdated bool
	if request.Code != "" && request.Code != quote.Code {
		bqbQuery.Space(`code = ?,`, request.Code)
	}
	if request.Date.IsZero() && request.Date.Equal(quote.Date) {
		bqbQuery.Space(`date = ?,`, request.Date)
	}
	if request.ExpiryDate.IsZero() && request.ExpiryDate.Equal(quote.ExpiryDate) && request.ExpiryDate.After(quote.Date) {
		bqbQuery.Space(`expiry_date = ?,`, request.ExpiryDate)
	}
	if request.Note != "" && request.Note != quote.Note {
		// -- Handle no note
		if request.Note == "no" {
			request.Note = ""
		}
		bqbQuery.Space(`note = ?,`, request.Note)
	}
	if request.Discount != 0 && request.Discount != quote.Discount {
		// -- Handle no discount
		if request.Discount == -1 {
			request.Discount = 0
		}
		bqbQuery.Space(`discount = ?,`, request.Discount)
		isDiscountUpdated = true
	}
	if request.ClientId != 0 && request.ClientId != quote.ClientId {
		bqbQuery.Space(`client_id = ?,`, request.ClientId)
	}
	if request.AccountId != 0 && request.AccountId != quote.AccountId {
		bqbQuery.Space(`account_id = ?,`, request.AccountId)
	}

	// -- Check if there are changes
	if len(bqbQuery.Parts) > 1 {
		// -- Remove last comma
		query, params, err = bqbQuery.Space(`status = ?, mid = ?, mtime = ? WHERE id = ?`, "Draft", tmpQuote.MId, tmpQuote.MTime, quoteId).ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing sql query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Update quote
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error updating quote: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Check if there are items to update, insert or delete
	if len(request.Items) == 0 && len(request.DeleteItems) == 0 {
		// Check if discount is updated
		if isDiscountUpdated {
			// -- Prepare sql query
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
		}

		tx.Commit()
		c.JSON(200, utils.NewResponse(200, fmt.Sprintf("quote %d updated successfully", quoteId), nil))
		return
	}

	// -- Separate quote items to update and insert
	var updateQuoteItems []models.QuoteItem
	var insertQuoteItems []models.QuoteItem
	for _, item := range request.Items {
		if item.Id == 0 {
			insertQuoteItems = append(insertQuoteItems, models.QuoteItem{
				Name:        item.Name,
				Description: item.Description,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
			})
		} else {
			updateQuoteItems = append(updateQuoteItems, models.QuoteItem{
				Id:          item.Id,
				Name:        item.Name,
				Description: item.Description,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
			})
		}
	}

	// -- Update quote items
	for _, item := range updateQuoteItems {
		if item.Name == "" && item.Description == "" && item.Quantity < 1 && item.UnitPrice <= 0 {
			continue
		}
		// -- Prepare sql query
		bqbQuery = bqb.New(`UPDATE "quote_item" SET`)
		if item.Name != "" {
			bqbQuery.Space(`name = ?`, item.Name)
		}
		if item.Description != "" {
			bqbQuery.Space(`description = ?`, item.Description)
		}
		if item.Quantity > 0 {
			bqbQuery.Space(`quantity = ?`, item.Quantity)
		}
		if item.UnitPrice > 0 {
			bqbQuery.Space(`unit_price = ?`, item.UnitPrice)
		}
		bqbQuery.Space(`mid, mtime = ? WHERE id = ?`, userId, quote.MTime, item.Id)
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error updating quote items: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Prepare sql query
	var createCount int
	bqbQuery = bqb.New(`INSERT INTO "quote_item" (quote_id, name, description, quantity, unit_price, cid, ctime, mid, mtime) VALUES`)
	for _, item := range insertQuoteItems {
		if item.Name == "" && item.Description == "" && item.Quantity < 1 && item.UnitPrice <= 0 {
			continue
		}
		bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?, ?),`, quoteId, item.Name, item.Description, item.Quantity, item.UnitPrice, tmpQuote.CId, tmpQuote.CTime, tmpQuote.MId, tmpQuote.MTime)
		createCount++
	}

	if createCount > 0 {
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
			log.Printf("Error inserting quote items: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	if len(request.DeleteItems) > 0 {
		// -- Prepare sql query
		query, params, err := bqb.New(`SELECT Count(*) FROM "quote_item" WHERE quote_id = ?`, quoteId).ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing sql query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		// -- Count quote items
		var count uint
		if err := tx.QueryRow(query, params...).Scan(&count); err != nil {
			tx.Rollback()
			log.Printf("Error counting quote items: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		if uint(len(request.DeleteItems))+1 > count {
			tx.Rollback()
			c.JSON(400, utils.NewErrorResponse(400, "invalid request body. delete items should not exceed total quote items or empty the quote items"))
			return
		}

		// -- Prepare sql query
		bqbQuery = bqb.New(`DELETE FROM "quote_item" WHERE id IN (`)
		for _, id := range request.DeleteItems {
			bqbQuery.Space(`?,`, id)
		}

		query, params, err = bqbQuery.ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing sql query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
		// -- Remove last comma and add closing bracket
		query = query[:len(query)-1] + ")"

		// -- Delete quote items
		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error deleting quote items: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Update quote total
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

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("quote %d updated successfully", quoteId), nil))
}

func (handler *QuoteHandler) Delete(c *gin.Context) {
	// -- Get id
	quoteId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. quote Id should be an integer"))
		return
	}

	// -- Begin transaction
	tx, err := handler.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT expiry_date, status FROM "quote" WHERE id = ?`, quoteId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get quote status
	var targetQuote models.Quote
	if err := tx.QueryRow(query, params...).Scan(&targetQuote.ExpiryDate, &targetQuote.Status); err != nil {
		tx.Rollback()

		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, "quote not found"))
			return
		}

		log.Printf("Error counting quote: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}
	if targetQuote.Status == "Accept" || targetQuote.Status == "Reject" || isExpired(targetQuote.ExpiryDate) {
		tx.Rollback()

		if isExpired(targetQuote.ExpiryDate) {
			c.JSON(400, utils.NewErrorResponse(400, "quote status is already expired. create a new quote instead"))
			return
		}

		c.JSON(400, utils.NewErrorResponse(400, "quote status is already "+strings.ToLower(targetQuote.Status)))
		return
	}

	// -- Prepare sql query
	query, params, err = bqb.New(`DELETE FROM "quote_item" WHERE quote_id = ?`, quoteId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing sql query: %v\n", err)
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

func isExpired(date time.Time) bool {
	return date.Before(time.Now())
}
