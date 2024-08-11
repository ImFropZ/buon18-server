package routes

import (
	"server/controllers"
	"server/database"
	"server/middlewares"
	"server/models/accounting"
	services "server/services/accounting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Accounting(e *gin.Engine, connection *database.Connection) {
	handler := controllers.AccountingHandler{
		DB:                            connection.DB,
		AccountingPaymentTermService:  &services.AccountingPaymentTermService{DB: connection.DB},
		AccountingAccountService:      &services.AccountingAccountService{DB: connection.DB},
		AccountingJournalService:      &services.AccountingJournalService{DB: connection.DB},
		AccountingJournalEntryService: &services.AccountingJournalEntryService{DB: connection.DB},
	}

	e.GET(
		"/api/accounting/accounts",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.VIEW}),
		middlewares.ValkeyCache[[]accounting.AccountingAccount](connection, "accounts"),
		handler.Accounts,
	)
	e.GET(
		"/api/accounting/accounts/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.VIEW}),
		handler.Account,
	)
	e.GET(
		"/api/accounting/payment-terms",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.VIEW}),
		middlewares.ValkeyCache[[]accounting.AccountingPaymentTerm](connection, "payment-terms"),
		handler.PaymentTerms,
	)
	e.GET(
		"/api/accounting/payment-terms/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.VIEW}),
		handler.PaymentTerm,
	)
	e.GET(
		"/api/accounting/journals",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.VIEW}),
		middlewares.ValkeyCache[[]accounting.AccountingJournal](connection, "journals"),
		handler.Journals,
	)
	e.GET(
		"/api/accounting/journals/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.VIEW}),
		handler.Journal,
	)
	e.GET(
		"/api/accounting/journal-entries",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.VIEW}),
		middlewares.ValkeyCache[[]accounting.AccountingJournalEntry](connection, "journal-entries"),
		handler.JournalEntries,
	)
}
