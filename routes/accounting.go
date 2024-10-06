package routes

import (
	"system.buon18.com/m/controllers"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/models/accounting"
	"system.buon18.com/m/services"
	accountingServices "system.buon18.com/m/services/accounting"
	"system.buon18.com/m/utils"

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
	e.DELETE(
		"/api/accounting/accounts/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.DELETE}),
		handler.DeleteAccount,
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
	e.DELETE(
		"/api/accounting/payment-terms/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.DELETE}),
		handler.DeletePaymentTerm,
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
	e.DELETE(
		"/api/accounting/journals/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.DELETE}),
		handler.DeleteJournal,
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
	e.DELETE(
		"/api/accounting/journal-entries/:id",
		middlewares.Authorize([]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.DELETE}),
		handler.DeleteJournalEntry,
	)
}
