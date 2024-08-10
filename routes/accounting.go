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
		AccountingAccountService:     &services.AccountingAccountService{DB: db},
		AccountingJournalService:     &services.AccountingJournalService{DB: db},
	}

	e.GET("/api/accounting/accounts", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.VIEW}), handler.Accounts)
	e.GET("/api/accounting/accounts/:id", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.VIEW}), handler.Account)
	e.GET("/api/accounting/payment-terms", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.VIEW}), handler.PaymentTerms)
	e.GET("/api/accounting/payment-terms/:id", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.VIEW}), handler.PaymentTerm)
	e.GET("/api/accounting/journals", middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.VIEW}), handler.Journals)
}
