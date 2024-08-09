package routes

import (
	"database/sql"
	"server/controllers"
	"server/middlewares"
	services "server/services/accounting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Accounting(e *gin.Engine, db *sql.DB) {
	handler := controllers.AccountingHandler{
		DB:                           db,
		AccountingPaymentTermService: &services.AccountingPaymentTermService{DB: db},
	}

	e.GET("/api/accounting/payment-terms", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.VIEW}), handler.PaymentTerms)
}
