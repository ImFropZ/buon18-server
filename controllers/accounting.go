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
	qp := utils.NewQueryParams()
	for _, filter := range accounting.AccountingPaymentTermAllowFilterFieldsAndOps {
		if validFilter, ok := c.GetQuery(filter); ok {
			qp.AddFilter(fmt.Sprintf(`"accounting.payment_term".%s=%s`, filter, validFilter))
		}
	}
	for _, sort := range accounting.AccountingPaymentTermAllowSortFields {
		if validSort, ok := c.GetQuery(fmt.Sprintf("sort-%s", sort)); ok {
			qp.AddOrderBy(fmt.Sprintf(`LOWER("limited_payment_terms".%s) %s`, sort, validSort))
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
