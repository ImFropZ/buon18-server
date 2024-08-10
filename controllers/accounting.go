package controllers

import (
	"database/sql"
	"fmt"
	"server/models/accounting"
	services "server/services/accounting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AccountingHandler struct {
	DB                           *sql.DB
	AccountingPaymentTermService *services.AccountingPaymentTermService
}

func (handler *AccountingHandler) PaymentTerms(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingPaymentTermAllowFilterFieldsAndOps, `"accounting.payment_term"`).
		PrepareSorts(c, accounting.AccountingPaymentTermAllowSortFields, `"limited_payment_terms"`).
		PreparePagination(c)

	paymentTerms, total, statusCode, err := handler.AccountingPaymentTermService.PaymentTerms(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"payment_terms": paymentTerms,
	}))
}

func (handler *AccountingHandler) PaymentTerm(c *gin.Context) {
	id := c.Param("id")

	paymentTerm, statusCode, err := handler.AccountingPaymentTermService.PaymentTerm(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"payment_term": paymentTerm,
	}))
}
