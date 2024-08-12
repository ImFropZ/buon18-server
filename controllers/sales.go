package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"server/models/sales"
	"server/services"
	"server/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type SalesHandler struct {
	DB            *sql.DB
	ServiceFacade *services.ServiceFacade
}

func (handler *SalesHandler) Quotations(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, sales.SalesQuotationAllowFilterFieldsAndOps, `"sales.quotation"`).
		PrepareSorts(c, sales.SalesQuotationAllowSortFields, `"limited_quotations"`).
		PreparePagination(c)

	quotations, total, statusCode, err := handler.ServiceFacade.SalesQuotationService.Quotations(qp)
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

	quotation, statusCode, err := handler.ServiceFacade.SalesQuotationService.Quotation(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"quotation": quotation,
	}))
}

func (handler *SalesHandler) CreateQuotation(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, "internal server error"))
		return
	}

	var quotation sales.SalesQuotationCreateRequest
	if err := c.ShouldBindJSON(&quotation); err != nil {
		log.Printf("Error binding JSON: %s", err)
		if strings.HasPrefix(err.Error(), "parsing time") {
			c.JSON(400, utils.NewErrorResponse(400, "invalid date format"))
			return
		}
		c.JSON(400, utils.NewErrorResponse(400, err.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(quotation); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.SalesQuotationService.CreateQuotation(&ctx, &quotation)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "quotation created successfully", nil))
}

func (handler *SalesHandler) Orders(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, sales.SalesOrderAllowFilterFieldsAndOps, `"sales.order"`).
		PrepareSorts(c, sales.SalesOrderAllowSortFields, `"limited_orders"`).
		PreparePagination(c)

	orders, total, statusCode, err := handler.ServiceFacade.SalesOrderService.Orders(qp)
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

	order, statusCode, err := handler.ServiceFacade.SalesOrderService.Order(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"order": order,
	}))
}
