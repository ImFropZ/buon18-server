package routes

import (
	"server/controllers"
	"server/database"
	"server/middlewares"
	"server/models/accounting"
	"server/services"
	accountingServices "server/services/accounting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

func Accounting(e *gin.Engine, connection *database.Connection) {
	handler := controllers.AccountingHandler{
		DB: connection.DB,
		ServiceFacade: &services.ServiceFacade{
			AccountingAccountService:      &accountingServices.AccountingAccountService{DB: connection.DB},
			AccountingPaymentTermService:  &accountingServices.AccountingPaymentTermService{DB: connection.DB},
			AccountingJournalService:      &accountingServices.AccountingJournalService{DB: connection.DB},
			AccountingJournalEntryService: &accountingServices.AccountingJournalEntryService{DB: connection.DB},
		},
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
	e.POST(
		"/api/accounting/accounts",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.CREATE}),
		handler.CreateAccount,
	)
	e.PATCH(
		"/api/accounting/accounts/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.UPDATE}),
		handler.UpdateAccount,
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
	e.POST(
		"/api/accounting/payment-terms",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.CREATE}),
		handler.CreatePaymentTerm,
	)
	e.PATCH(
		"/api/accounting/payment-terms/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.UPDATE}),
		handler.UpdatePaymentTerm,
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
	e.POST(
		"/api/accounting/journals",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.CREATE}),
		handler.CreateJournal,
	)
	e.PATCH(
		"/api/accounting/journals/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.UPDATE}),
		handler.UpdateJournal,
	)
	e.GET(
		"/api/accounting/journal-entries",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.VIEW}),
		middlewares.ValkeyCache[[]accounting.AccountingJournalEntry](connection, "journal-entries"),
		handler.JournalEntries,
	)
	e.GET(
		"/api/accounting/journal-entries/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.VIEW}),
		handler.JournalEntry,
	)
	e.POST(
		"/api/accounting/journal-entries",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.CREATE}),
		handler.CreateJournalEntry,
	)
	e.PATCH(
		"/api/accounting/journal-entries/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.UPDATE}),
		handler.UpdateJournalEntry,
	)
}
