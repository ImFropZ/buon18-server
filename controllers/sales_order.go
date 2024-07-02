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

	c.JSON(200, utils.NewResponse(200, "success", salesOrder.ToResponse()))
}

func (handler *SalesOrderHandler) List(c *gin.Context) {
	paginationQueryParams := utils.PaginationQueryParams{
		Offset: 0,
		Limit:  10,
	}

	// -- Parse query params
	paginationQueryParams.Parse(c)

	// -- Prepare sql query
	query, params, err := bqb.New(`SELECT id, code, COALESCE(note, ''), status, accept_date, delivery_date, quote_id, cid 
	FROM 
		"sales_order" 
	ORDER BY id
	LIMIT ? OFFSET ?`, paginationQueryParams.Limit, paginationQueryParams.Offset).ToPgsql()
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

	c.JSON(200, utils.NewResponse(200, "success", models.SalesOrdersToResponse(salesOrders)))
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
