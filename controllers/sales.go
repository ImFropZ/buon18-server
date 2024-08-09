package controllers

import (
	"database/sql"
	"fmt"
	"server/models/sales"
	services "server/services/sales"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type SalesHandler struct {
	DB                    *sql.DB
	SalesQuotationService *services.SalesQuotationService
}

func (handler *SalesHandler) Quotations(c *gin.Context) {
	qp := utils.NewQueryParams()
	for _, filter := range sales.SalesQuotationAllowFilterFieldsAndOps {
		if validFilter, ok := c.GetQuery(filter); ok {
			qp.AddFilter(fmt.Sprintf(`"sales.quotation".%s=%s`, filter, validFilter))
		}
	}
	for _, sort := range sales.SalesQuotationAllowSortFields {
		if validSort, ok := c.GetQuery(fmt.Sprintf("sort-%s", sort)); ok {
			qp.AddOrderBy(fmt.Sprintf(`LOWER("limited_quotations".%s) %s`, sort, validSort))
		}
	}
	for _, pagination := range []string{"offset", "limit"} {
		if validPagination, ok := c.GetQuery(pagination); ok {
			if pagination == "offset" {
				qp.AddOffset(utils.StrToInt(validPagination, 0))
			} else {
				qp.AddLimit(utils.StrToInt(validPagination, 10))
			}
		}
	}

	quotations, total, statusCode, err := handler.SalesQuotationService.Quotations(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"quotations": quotations,
	}))
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
