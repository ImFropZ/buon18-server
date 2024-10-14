package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"system.buon18.com/m/controllers"
	"system.buon18.com/m/database"
	"system.buon18.com/m/middlewares"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"

	accountingServices "system.buon18.com/m/services/accounting"
)

func AccountingRoutes(r *mux.Router, con *database.Connection) {
	controller := controllers.AccountingHandler{
		DB: con.DB,
		ServiceFacade: &services.ServiceFacade{
			AccountingAccountService:      &accountingServices.AccountingAccountService{DB: con.DB},
			AccountingPaymentTermService:  &accountingServices.AccountingPaymentTermService{DB: con.DB},
			AccountingJournalService:      &accountingServices.AccountingJournalService{DB: con.DB},
			AccountingJournalEntryService: &accountingServices.AccountingJournalEntryService{DB: con.DB},
		},
	}

	r.Handle(
		"/accounts",
		middlewares.Authorize(
			http.HandlerFunc(controller.Accounts),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/accounts/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.Account),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/accounts",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateAccount),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/accounts/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateAccount),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/accounts",
		middlewares.Authorize(
			http.HandlerFunc(controller.DeleteAccounts),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_ACCOUNTS.DELETE},
		),
	).Methods("DELETE", "OPTIONS")

	r.Handle(
		"/payment-terms",
		middlewares.Authorize(
			http.HandlerFunc(controller.PaymentTerms),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/payment-terms/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.PaymentTerm),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/payment-terms",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreatePaymentTerm),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/payment-terms/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdatePaymentTerm),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/payment-terms",
		middlewares.Authorize(
			http.HandlerFunc(controller.DeletePaymentTerms),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_PAYMENT_TERMS.DELETE},
		),
	).Methods("DELETE", "OPTIONS")

	r.Handle(
		"/journals",
		middlewares.Authorize(
			http.HandlerFunc(controller.Journals),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/journals/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.Journal),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/journals",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateJournal),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/journals/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateJournal),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/journals",
		middlewares.Authorize(
			http.HandlerFunc(controller.DeleteJournals),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNALS.DELETE},
		),
	).Methods("DELETE", "OPTIONS")

	r.Handle(
		"/journal-entries",
		middlewares.Authorize(
			http.HandlerFunc(controller.JournalEntries),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/journal-entries/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.JournalEntry),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.VIEW},
		),
	).Methods("GET", "OPTIONS")

	r.Handle(
		"/journal-entries",
		middlewares.Authorize(
			http.HandlerFunc(controller.CreateJournalEntry),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.CREATE},
		),
	).Methods("POST", "OPTIONS")

	r.Handle(
		"/journal-entries/{id:[0-9]+}",
		middlewares.Authorize(
			http.HandlerFunc(controller.UpdateJournalEntry),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.UPDATE},
		),
	).Methods("PATCH", "OPTIONS")

	r.Handle(
		"/journal-entries",
		middlewares.Authorize(
			http.HandlerFunc(controller.DeleteJournalEntries),
			[]string{utils.PREDEFINED_PERMISSIONS.FULL_ACCOUNTING, utils.PREDEFINED_PERMISSIONS.ACCOUNTING_JOURNAL_ENTRIES.DELETE},
		),
	).Methods("DELETE", "OPTIONS")
}
