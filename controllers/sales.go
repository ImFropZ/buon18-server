package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"server/models/sales"
	services "server/services/sales"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type SalesHandler struct {
	DB                    *sql.DB
	SalesQuotationService *services.SalesQuotationService
	SalesOrderService     *services.SalesOrderService
}

func (handler *SalesHandler) Quotations(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, sales.SalesQuotationAllowFilterFieldsAndOps, `"sales.quotation"`).
		PrepareSorts(c, sales.SalesQuotationAllowSortFields, `"limited_quotations"`).
		PreparePagination(c)

	quotations, total, statusCode, err := handler.SalesQuotationService.Quotations(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"quotations": quotations,
	}))

	c.Set("total", total)
	if quotationsByte, err := json.Marshal(quotations); err == nil {
		c.Set("response", quotationsByte)
	}
}

func (handler *SalesHandler) Quotation(c *gin.Context) {
	id := c.Param("id")

	quotation, statusCode, err := handler.SalesQuotationService.Quotation(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"quotation": quotation,
	}))
}

func (handler *SalesHandler) Orders(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, sales.SalesOrderAllowFilterFieldsAndOps, `"sales.order"`).
		PrepareSorts(c, sales.SalesOrderAllowSortFields, `"limited_orders"`).
		PreparePagination(c)

	orders, total, statusCode, err := handler.SalesOrderService.Orders(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"orders": orders,
	}))

	c.Set("total", total)
	if ordersByte, err := json.Marshal(orders); err == nil {
		c.Set("response", ordersByte)
	}
}

func (handler *SalesHandler) Order(c *gin.Context) {
	id := c.Param("id")

	order, statusCode, err := handler.SalesOrderService.Order(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"order": order,
	}))
}
