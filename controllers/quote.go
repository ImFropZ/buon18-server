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

func prepareQuoteQuery(c *gin.Context, bqbQuery *bqb.Query) {
	// -- Apply query params
	bqbQuery.Space("WHERE")
	if str, ok := c.GetQuery("code_ilike"); ok {
		bqbQuery.Space(`"quote".code ILIKE ? AND`, "%"+str+"%")
	}
	if str, ok := c.GetQuery("client_id_ilike"); ok {
		bqbQuery.Space(`"quote".client_id ILIKE ? AND`, "%"+str+"%")
	}
	if str, ok := c.GetQuery("account_id_ilike"); ok {
		bqbQuery.Space(`"quote".account_id ILIKE ? AND`, "%"+str+"%")
	}
	if str, ok := c.GetQuery("created_by_ilike"); ok {
		bqbQuery.Space(`"quote".cid ILIKE ? AND`, "%"+str+"%")
	}
	if str, ok := c.GetQuery("date_min"); ok {
		bqbQuery.Space(`"quote".date >= ? AND`, str)
	}
	if str, ok := c.GetQuery("date_max"); ok {
		bqbQuery.Space(`"quote".date <= ? AND`, str)
	}
	if str, ok := c.GetQuery("expiry_date_min"); ok {
		bqbQuery.Space(`"quote".expiry_date >= ? AND`, str)
	}
	if str, ok := c.GetQuery("expiry_date_max"); ok {
		bqbQuery.Space(`"quote".expiry_date <= ? AND`, str)
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

type CreateQuoteRequest struct {
	Code       string    `json:"code" binding:"required"`
	Date       time.Time `json:"date" binding:"required"`
	ExpiryDate time.Time `json:"expiry_date" binding:"required"`
	Note       *string   `json:"note"`
	Discount   *float64  `json:"discount"`
	ClientId   uint      `json:"client_id" binding:"required"`
	AccountId  uint      `json:"account_id" binding:"required"`
	Items      []struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
		Quantity    uint    `json:"quantity" binding:"required"`
		UnitPrice   float64 `json:"unit_price" binding:"required"`
	} `json:"items" binding:"required"`
}

type UpdateQuoteRequest struct {
	Code       *string    `json:"code"`
	Date       *time.Time `json:"date"`
	ExpiryDate *time.Time `json:"expiry_date"`
	Note       *string    `json:"note"`
	Discount   *float64   `json:"discount"`
	ClientId   *uint      `json:"client_id"`
	AccountId  *uint      `json:"account_id"`
	Items      *[]struct {
		Id          *uint    `json:"id"`
		Name        *string  `json:"name"`
		Description *string  `json:"description"`
		Quantity    *uint    `json:"quantity"`
		UnitPrice   *float64 `json:"unit_price"`
	} `json:"items"`

	// -- Delete items
	DeleteItems *[]uint `json:"delete_item_ids"`
}

type CreateSalesOrderRequest struct {
	Code         string     `json:"code" binding:"required"`
	AcceptDate   *time.Time `json:"accept_date"`
	DeliveryDate time.Time  `json:"delivery_date" binding:"required"`
	Note         *string    `json:"note"`
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
	query, params, err := bqb.New(`SELECT "quote".id, "quote".code, "quote".date, "quote".expiry_date, COALESCE("quote".note, ''), "quote".subtotal, "quote".discount, "quote".total, "quote".client_id, "quote".account_id, "quote".status, "quote".cid, COALESCE("quote_item".id, 0), COALESCE("quote_item".name, ''), COALESCE("quote_item".description, ''), COALESCE("quote_item".quantity, 0), COALESCE("quote_item".unit_price, 0)
	FROM "quote"
	LEFT JOIN "quote_item" ON "quote_item".quote_id = "quote".id
	WHERE "quote".id = ? 
	ORDER BY "quote_item".id`, quoteId).ToPgsql()
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
	bqbQuery := bqb.New(`SELECT "quote".id, "quote".code, "quote".date, "quote".expiry_date, COALESCE("quote".note, ''), "quote".subtotal, "quote".discount, "quote".total, "quote".client_id, "quote".account_id, "quote".status, "quote".cid, COALESCE("quote_item".id, 0), COALESCE("quote_item".name, ''), COALESCE("quote_item".description, ''), COALESCE("quote_item".quantity, 0), COALESCE("quote_item".unit_price, 0)
	FROM "quote"
	LEFT JOIN "quote_item" ON "quote_item".quote_id = "quote".id`)

	// -- Apply query params
	prepareQuoteQuery(c, bqbQuery)

	// -- Complete query
	bqbQuery.Space(`ORDER BY "quote".id, "quote_item".id OFFSET ? LIMIT ?`, paginationQueryParams.Offset, paginationQueryParams.Limit)

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

	prepareQuoteQuery(c, bqbQuery)

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
	total := subtotal
	if request.Discount != nil { // WARN: Discount nil means 0 discount
		total -= *request.Discount
	}

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

	var acceptDate time.Time
	if request.AcceptDate == nil {
		acceptDate = time.Now()
	} else {
		acceptDate = *request.AcceptDate
	}

	// -- Prepare sql query
	query, params, err = bqb.New(`INSERT INTO "sales_order" 
	(code, accept_date, delivery_date, note, quote_id, cid, ctime, mid, mtime) 
	VALUES 
	(?, ?, ?, ?, ?, ?, ?, ?, ?)`, request.Code, acceptDate, request.DeliveryDate, *request.Note, quoteId, tmpSalesOrder.CId, tmpSalesOrder.CTime, tmpSalesOrder.MId, tmpSalesOrder.MTime).ToPgsql()
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
	var req UpdateQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body."))
		return
	}

	// -- Check if all fields are nil
	if utils.IsAllFieldsNil(&req) {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body. at least one field should be updated"))
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
	if req.Code != nil {
		bqbQuery.Space(`code = ?,`, *req.Code)
	}
	if req.Date != nil || req.ExpiryDate != nil {
		if req.Date != nil && req.ExpiryDate != nil {
			if !req.Date.Before(*req.ExpiryDate) {
				bqbQuery.Space(`date = ?, expiry_date = ?,`, *req.Date, *req.ExpiryDate)
			} else {
				c.JSON(400, utils.NewErrorResponse(400, "invalid request body. expiry_date should be after date"))
				return
			}
		} else if req.Date != nil {
			if !req.Date.Before(quote.ExpiryDate) {
				bqbQuery.Space(`date = ?,`, *req.Date)
			} else {
				c.JSON(400, utils.NewErrorResponse(400, "invalid request body. expiry_date should be after date"))
				return
			}
		} else if req.ExpiryDate != nil {
			if !quote.Date.Before(*req.ExpiryDate) {
				bqbQuery.Space(`expiry_date = ?,`, *req.ExpiryDate)
			} else {
				c.JSON(400, utils.NewErrorResponse(400, "invalid request body. expiry_date should be after date"))
				return
			}
		}
	}

	if req.Note != nil {
		bqbQuery.Space(`note = ?,`, req.Note)
	}
	if req.Discount != nil {
		bqbQuery.Space(`discount = ?,`, req.Discount)
		isDiscountUpdated = true
	}
	if req.ClientId != nil {
		bqbQuery.Space(`client_id = ?,`, req.ClientId)
	}
	if req.AccountId != nil {
		bqbQuery.Space(`account_id = ?,`, req.AccountId)
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
	if (req.Items == nil || len(*req.Items) == 0) && (req.DeleteItems == nil || len(*req.DeleteItems) == 0) {
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
	var insertQuoteItems []models.QuoteItemCreate
	var updateQuoteItems []models.QuoteItemUpdate
	for _, item := range *req.Items {
		if item.Id == nil {
			if item.Name == nil || item.Quantity == nil || item.UnitPrice == nil {
				continue
			}
			insertQuoteItems = append(insertQuoteItems, models.QuoteItemCreate{
				Name:        *item.Name,
				Description: item.Description,
				Quantity:    *item.Quantity,
				UnitPrice:   *item.UnitPrice,
			})
		} else {
			updateQuoteItems = append(updateQuoteItems, models.QuoteItemUpdate{
				Id:          *item.Id,
				Name:        item.Name,
				Description: item.Description,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
			})
		}
	}

	// -- Update quote items
	for _, item := range updateQuoteItems {
		if item.Name == nil && item.Description == nil && item.Quantity == nil && item.UnitPrice == nil {
			continue
		}
		// -- Prepare sql query
		bqbQuery = bqb.New(`UPDATE "quote_item" SET`)
		if item.Name != nil {
			bqbQuery.Space(`name = ?,`, *item.Name)
		}
		if item.Description != nil {
			bqbQuery.Space(`description = ?,`, *item.Description)
		}
		if item.Quantity != nil {
			bqbQuery.Space(`quantity = ?,`, *item.Quantity)
		}
		if item.UnitPrice != nil {
			bqbQuery.Space(`unit_price = ?,`, *item.UnitPrice)
		}
		query, params, err = bqbQuery.Space(`mid = ?, mtime = ? WHERE id = ?`, userId, quote.MTime, item.Id).ToPgsql()
		if err != nil {
			tx.Rollback()
			log.Printf("Error preparing sql query: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}

		if _, err := tx.Exec(query, params...); err != nil {
			tx.Rollback()
			log.Printf("Error updating quote items: %v\n", err)
			c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
			return
		}
	}

	// -- Prepare sql query
	hasCreate := false
	bqbQuery = bqb.New(`INSERT INTO "quote_item" (quote_id, name, description, quantity, unit_price, cid, ctime, mid, mtime) VALUES`)
	for _, item := range insertQuoteItems {
		if item.Description == nil && item.Quantity < 1 && item.UnitPrice <= 0 {
			continue
		}
		bqbQuery.Space(`(?, ?, ?, ?, ?, ?, ?, ?, ?),`, quoteId, item.Name, item.Description, item.Quantity, item.UnitPrice, tmpQuote.CId, tmpQuote.CTime, tmpQuote.MId, tmpQuote.MTime)
		hasCreate = true
	}

	if hasCreate {
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

	if req.DeleteItems != nil && len(*req.DeleteItems) > 0 {
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

		if uint(len(*req.DeleteItems))+1 > count {
			tx.Rollback()
			c.JSON(400, utils.NewErrorResponse(400, "invalid request body. delete items should not exceed total quote items or empty the quote items"))
			return
		}

		// -- Prepare sql query
		bqbQuery = bqb.New(`DELETE FROM "quote_item" WHERE id IN (`)
		for _, id := range *req.DeleteItems {
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
