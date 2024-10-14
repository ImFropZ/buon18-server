package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"system.buon18.com/m/models"
	"system.buon18.com/m/models/accounting"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"
)

type AccountingHandler struct {
	DB            *sql.DB
	ServiceFacade *services.ServiceFacade
}

func (handler *AccountingHandler) Accounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(accounting.AccountingAccount{}, r, `"accounting.account"`).
		PrepareSorts(accounting.AccountingAccount{}, r, `"accounting.account"`).
		PrepareLimitAndOffset(r)

	accounts, total, statusCode, err := handler.ServiceFacade.AccountingAccountService.Accounts(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"accounts": accounts,
	}))
}

func (handler *AccountingHandler) Account(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	account, statusCode, err := handler.ServiceFacade.AccountingAccountService.Account(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"account": account,
	}))
}

func (handler *AccountingHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingAccountCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingAccountService.CreateAccount(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *AccountingHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingAccountUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]
	statusCode, err := handler.ServiceFacade.AccountingAccountService.UpdateAccount(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *AccountingHandler) DeleteAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingAccountService.DeleteAccounts(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}

func (handler *AccountingHandler) PaymentTerms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(accounting.AccountingPaymentTerm{}, r, `"accounting.payment_term"`).
		PrepareSorts(accounting.AccountingPaymentTerm{}, r, `"limited_payment_terms"`).
		PrepareLimitAndOffset(r)

	paymentTerms, total, statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.PaymentTerms(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"payment_terms": paymentTerms,
	}))
}

func (handler *AccountingHandler) PaymentTerm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	paymentTerm, statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.PaymentTerm(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"payment_term": paymentTerm,
	}))
}

func (handler *AccountingHandler) CreatePaymentTerm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingPaymentTermCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.CreatePaymentTerm(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *AccountingHandler) UpdatePaymentTerm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingPaymentTermUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]
	statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.UpdatePaymentTerm(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *AccountingHandler) DeletePaymentTerms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.DeletePaymentTerms(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}

func (handler *AccountingHandler) Journals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(accounting.AccountingJournal{}, r, `"accounting.journal"`).
		PrepareSorts(accounting.AccountingJournal{}, r, `"limited_journals"`).
		PrepareLimitAndOffset(r)

	journals, total, statusCode, err := handler.ServiceFacade.AccountingJournalService.Journals(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"journals": journals,
	}))
}

func (handler *AccountingHandler) Journal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	journal, statusCode, err := handler.ServiceFacade.AccountingJournalService.Journal(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"journal": journal,
	}))
}

func (handler *AccountingHandler) CreateJournal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingJournalCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalService.CreateJournal(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *AccountingHandler) UpdateJournal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingJournalUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]
	statusCode, err := handler.ServiceFacade.AccountingJournalService.UpdateJournal(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *AccountingHandler) DeleteJournals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalService.DeleteJournals(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}

func (handler *AccountingHandler) JournalEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qp := utils.NewQueryParams().
		PrepareFilters(accounting.AccountingJournalEntry{}, r, `"accounting.journal_entry"`).
		PrepareSorts(accounting.AccountingJournalEntry{}, r, `"limited_journal_entries"`).
		PrepareLimitAndOffset(r)

	journalEntries, total, statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.JournalEntries(qp)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.Header().Set("X-Total-Count", fmt.Sprintf("%d", total))
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"journal_entries": journalEntries,
	}))
}

func (handler *AccountingHandler) JournalEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	journalEntry, statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.JournalEntry(id)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "", map[string]interface{}{
		"journal_entry": journalEntry,
	}))
}

func (handler *AccountingHandler) CreateJournalEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingJournalEntryCreateRequest](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.CreateJournalEntry(ctx, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "created", nil))
}

func (handler *AccountingHandler) UpdateJournalEntry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value(utils.CtxKey{}).(*utils.CtxValue)
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[accounting.AccountingJournalEntryUpdateRequest](r, w, true)
	if !ok {
		return
	}

	id := mux.Vars(r)["id"]
	statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.UpdateJournalEntry(ctx, id, &req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "updated", nil))
}

func (handler *AccountingHandler) DeleteJournalEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// -- Parse request
	req, ok := utils.ValidateRequest[models.CommonDelete](r, w, false)
	if !ok {
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.DeleteJournalEntries(&req)
	if err != nil {
		msg, clientErr, code := utils.ServerToClientError(err)
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(utils.NewErrorResponse(code, msg, clientErr, nil))
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(utils.NewResponse(statusCode, "deleted", nil))
}
