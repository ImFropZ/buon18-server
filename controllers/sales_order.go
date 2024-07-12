package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"server/models"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type UpdateSalesOrderStatusRequest struct {
	Action string `json:"action" binding:"required"`
}

type SalesOrderHandler struct {
	DB *sql.DB
}

func (handler *SalesOrderHandler) First(c *gin.Context) {
	// -- Get id
	salesOrderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid sales order Id. sales order Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT id, code, COALESCE(note, ''), status, accept_date, delivery_date, quote_id, cid 
	FROM 
		"sales_order" 
	WHERE id = ?`, salesOrderId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query sales order
	var salesOrder models.SalesOrder
	if err := handler.DB.QueryRow(query, params...).Scan(&salesOrder.Id, &salesOrder.Code, &salesOrder.Note, &salesOrder.Status, &salesOrder.AcceptDate, &salesOrder.DeliveryDate, &salesOrder.QuoteId, &salesOrder.CId); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("sales order %d not found", salesOrderId)))
			return
		}

		log.Printf("Error querying sales order: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"sales_order": salesOrder.ToResponse(),
	}))
}

func (handler *SalesOrderHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Prepare sql query
	bqbQuery := bqb.New(`SELECT id, code, COALESCE(note, ''), status, accept_date, delivery_date, quote_id, cid 
	FROM 
		"sales_order"`)

	// -- Add query if exists
	if paginationQueryParams.Query != "" {
		bqbQuery.Space(`WHERE code ILIKE ?`, "%"+paginationQueryParams.Query+"%")
	}

	// -- Complete query
	bqbQuery.Space("ORDER BY id OFFSET ? LIMIT ?", paginationQueryParams.Offset, paginationQueryParams.Limit)

	query, params, err := bqbQuery.ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query sales orders
	var salesOrders []models.SalesOrder
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error querying sales orders: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		defer rows.Close()

		for rows.Next() {
			var so models.SalesOrder
			if err := rows.Scan(&so.Id, &so.Code, &so.Note, &so.Status, &so.AcceptDate, &so.DeliveryDate, &so.QuoteId, &so.CId); err != nil {
				log.Printf("Error scanning sales order: %v", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			salesOrders = append(salesOrders, so)
		}
	}

	// -- Count total sales orders
	bqbQuery = bqb.New(`SELECT COUNT(*) FROM "sales_order"`)

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
		log.Printf("Error getting total sales orders: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", gin.H{
		"total":        total,
		"sales_orders": models.SalesOrdersToResponse(salesOrders),
	}))
}

func (handler *SalesOrderHandler) CreateInvoice(c *gin.Context) {
	// -- Get id
	salesOrderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid sales order Id. sales order Id should be an integer"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT status FROM "sales_order" WHERE id = ?`, salesOrderId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get status
	var status string
	if err := handler.DB.QueryRow(query, params...).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("sales order %d not found", salesOrderId)))
			return
		}

		log.Printf("Error getting status: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	if status != "Done" {
		c.JSON(400, utils.NewErrorResponse(400, "invalid action. sales order is not done"))
		return
	}

	// -- Prepare sql query
	query, params, err = bqb.New(`SELECT id, code, COALESCE(note, ''), accept_date, delivery_date, quote_id, cid 
	FROM 
		"sales_order"
	WHERE id = ?`, salesOrderId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query sales order
	var salesOrder models.SalesOrder
	if err := handler.DB.QueryRow(query, params...).Scan(&salesOrder.Id, &salesOrder.Code, &salesOrder.Note, &salesOrder.AcceptDate, &salesOrder.DeliveryDate, &salesOrder.QuoteId, &salesOrder.CId); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("sales order %d not found", salesOrderId)))
			return
		}

		log.Printf("Error querying sales order: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query (GET USER)
	query, params, err = bqb.New(`SELECT name, email, role FROM "user" WHERE id = ?`, salesOrder.CId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query user
	var user models.User
	if err := handler.DB.QueryRow(query, params...).Scan(&user.Name, &user.Email, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("user %d not found", salesOrder.CId)))
			return
		}

		log.Printf("Error querying user: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query (GET Quote)
	query, params, err = bqb.New(`SELECT "quote".code, "quote".account_id, "quote".client_id, "quote".subtotal, "quote".discount, "quote".total, "quote_item".name, COALESCE("quote_item".description, ''), "quote_item".quantity, "quote_item".unit_price
	FROM "quote"
	LEFT JOIN "quote_item" ON "quote_item".quote_id = "quote".id
	WHERE "quote".id = ? 
	ORDER BY "quote_item".id`, salesOrder.QuoteId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query quote
	var quote models.Quote
	quote.QuoteItems = make([]models.QuoteItem, 0)
	if rows, err := handler.DB.Query(query, params...); err != nil {
		log.Printf("Error finding quote in database: %v\n", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		for rows.Next() {
			var scanQuoteItem models.QuoteItem
			if err := rows.Scan(&quote.Code, &quote.AccountId, &quote.ClientId, &quote.Subtotal, &quote.Discount, &quote.Total, &scanQuoteItem.Name, &scanQuoteItem.Description, &scanQuoteItem.Quantity, &scanQuoteItem.UnitPrice); err != nil {
				log.Printf("Error scanning quote from database: %v\n", err)
				c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
				return
			}

			quote.QuoteItems = append(quote.QuoteItems, scanQuoteItem)
		}
	}

	// -- Prepare sql query (GET Account)
	query, params, err = bqb.New(`SELECT 
	code, name, COALESCE(email, ''), COALESCE(address, ''), phone, COALESCE(secondary_phone, '')
	FROM
		"account"
	WHERE
		id = ?`, quote.AccountId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query account
	var account models.Account
	if err := handler.DB.QueryRow(query, params...).Scan(&account.Code, &account.Name, &account.Email, &account.Address, &account.Phone, &account.SecondaryPhone); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("account %d not found", user.Id)))
			return
		}

		log.Printf("Error querying account: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Prepare sql query (GET Client)
	query, params, err = bqb.New(`SELECT 
	code, name, COALESCE(address, ''), phone
	FROM "client" WHERE id = ?`, quote.ClientId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Query client
	var client models.Client
	if err := handler.DB.QueryRow(query, params...).Scan(&client.Code, &client.Name, &client.Address, &client.Phone); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("client %d not found", account.Id)))
			return
		}

		log.Printf("Error querying client: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, "success", utils.InvoiceResponse(salesOrder, quote, client, account, user)))
}

func (handler *SalesOrderHandler) UpdateStatus(c *gin.Context) {
	// -- Get user id
	var userId uint
	if id, ok := c.Get("user_id"); !ok {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	} else {
		userId = id.(uint)
	}

	// -- Get id
	salesOrderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid quote Id. sales order Id should be an integer"))
		return
	}

	// -- Parse request
	var request UpdateSalesOrderStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, "invalid request body. action is required"))
		return
	}

	// -- Validate action
	caser := cases.Title(language.English)
	if action := caser.String(request.Action); action == "On-Going" || action == "Sent" || action == "Done" || action == "Cancel" {
		request.Action = action
	} else {
		c.JSON(400, utils.NewErrorResponse(400, "invalid action. action should be one of on-going, sent, done, or cancel"))
		return
	}

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT status FROM "sales_order" WHERE id = ?`, salesOrderId).ToPgsql()
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Begin transaction
	tx, err := handler.DB.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Get status
	var status string
	err = tx.QueryRow(query, params...).Scan(&status)
	if err != nil {
		tx.Rollback()

		if err == sql.ErrNoRows {
			c.JSON(404, utils.NewErrorResponse(404, fmt.Sprintf("sales order %d not found", salesOrderId)))
			return
		}

		log.Printf("Error getting status: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	if status == request.Action {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, fmt.Sprintf("invalid action. sales order is already %s", status)))
		return
	}

	// -- Validate status
	if status == "Done" {
		tx.Rollback()
		c.JSON(400, utils.NewErrorResponse(400, "invalid action. sales order is already done"))
		return
	}

	tmpSalesOrder := models.SalesOrder{}
	tmpSalesOrder.PrepareForUpdate(userId)

	// -- Prepare sql query
	query, params, err = bqb.New(`UPDATE "sales_order" SET status = ?, mid = ?, mtime = ? WHERE id = ?`, request.Action, tmpSalesOrder.MId, tmpSalesOrder.MTime, salesOrderId).ToPgsql()
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing query: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Update status
	_, err = tx.Exec(query, params...)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating status: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	// -- Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	c.JSON(200, utils.NewResponse(200, fmt.Sprintf("sales order %d updated successfully", salesOrderId), nil))
}
