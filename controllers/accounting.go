package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"system.buon18.com/m/models/accounting"
	"system.buon18.com/m/services"
	"system.buon18.com/m/utils"

	"github.com/gin-gonic/gin"
)

type AccountingHandler struct {
	DB            *sql.DB
	ServiceFacade *services.ServiceFacade
}

func (handler *AccountingHandler) Accounts(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingAccountAllowFilterFieldsAndOps, `"accounting.account"`).
		PrepareSorts(c, accounting.AccountingAccountAllowSortFields, `"accounting.account"`).
		PreparePagination(c)

	accounts, total, statusCode, err := handler.ServiceFacade.AccountingAccountService.Accounts(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"accounts": accounts,
	}))

	c.Set("total", total)
	if accountsByte, err := json.Marshal(accounts); err == nil {
		c.Set("response", accountsByte)
	}
}

func (handler *AccountingHandler) Account(c *gin.Context) {
	id := c.Param("id")

	account, statusCode, err := handler.ServiceFacade.AccountingAccountService.Account(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"account": account,
	}))
}

func (handler *AccountingHandler) CreateAccount(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	var account accounting.AccountingAccountCreateRequest
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, utils.ErrInternalServer.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(account); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingAccountService.CreateAccount(&ctx, &account)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "account created successfully", nil))
}

func (handler *AccountingHandler) UpdateAccount(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	id := c.Param("id")

	var account accounting.AccountingAccountUpdateRequest
	if err := c.ShouldBindJSON(&account); err != nil {
		log.Printf("%v", err)
		c.JSON(400, utils.NewErrorResponse(400, utils.ErrInternalServer.Error()))
		return
	}

	if utils.IsAllFieldsNil(&account) {
		c.JSON(400, utils.NewErrorResponse(400, "no fields to update"))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(account); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingAccountService.UpdateAccount(&ctx, id, &account)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "account updated successfully", nil))
}

func (handler *AccountingHandler) DeleteAccount(c *gin.Context) {
	id := c.Param("id")

	statusCode, err := handler.ServiceFacade.AccountingAccountService.DeleteAccount(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "account deleted successfully", nil))
}

func (handler *AccountingHandler) PaymentTerms(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingPaymentTermAllowFilterFieldsAndOps, `"accounting.payment_term"`).
		PrepareSorts(c, accounting.AccountingPaymentTermAllowSortFields, `"limited_payment_terms"`).
		PreparePagination(c)

	paymentTerms, total, statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.PaymentTerms(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"payment_terms": paymentTerms,
	}))

	c.Set("total", total)
	if paymentTermsByte, err := json.Marshal(paymentTerms); err == nil {
		c.Set("response", paymentTermsByte)
	}
}

func (handler *AccountingHandler) PaymentTerm(c *gin.Context) {
	id := c.Param("id")

	paymentTerm, statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.PaymentTerm(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"payment_term": paymentTerm,
	}))
}

func (handler *AccountingHandler) CreatePaymentTerm(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	var paymentTerm accounting.AccountingPaymentTermCreateRequest
	if err := c.ShouldBindJSON(&paymentTerm); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, utils.ErrInternalServer.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(paymentTerm); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.CreatePaymentTerm(&ctx, &paymentTerm)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "payment term created successfully", nil))
}

func (handler *AccountingHandler) UpdatePaymentTerm(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	id := c.Param("id")

	var paymentTerm accounting.AccountingPaymentTermUpdateRequest
	if err := c.ShouldBindJSON(&paymentTerm); err != nil {
		log.Printf("%v", err)
		c.JSON(400, utils.NewErrorResponse(400, utils.ErrInternalServer.Error()))
		return
	}

	if utils.IsAllFieldsNil(&paymentTerm) {
		c.JSON(400, utils.NewErrorResponse(400, "no fields to update"))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(paymentTerm); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.UpdatePaymentTerm(&ctx, id, &paymentTerm)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "payment term updated successfully", nil))
}

func (handler *AccountingHandler) DeletePaymentTerm(c *gin.Context) {
	id := c.Param("id")

	statusCode, err := handler.ServiceFacade.AccountingPaymentTermService.DeletePaymentTerm(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "payment term deleted successfully", nil))
}

func (handler *AccountingHandler) Journals(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingJournalAllowFilterFieldsAndOps, `"accounting.journal"`).
		PrepareSorts(c, accounting.AccountingJournalAllowSortFields, `"limited_journals"`).
		PreparePagination(c)

	journals, total, statusCode, err := handler.ServiceFacade.AccountingJournalService.Journals(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"journals": journals,
	}))

	c.Set("total", total)
	if journalsByte, err := json.Marshal(journals); err == nil {
		c.Set("response", journalsByte)
	}
}

func (handler *AccountingHandler) Journal(c *gin.Context) {
	id := c.Param("id")

	journal, statusCode, err := handler.ServiceFacade.AccountingJournalService.Journal(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"journal": journal,
	}))
}

func (handler *AccountingHandler) CreateJournal(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	var journal accounting.AccountingJournalCreateRequest
	if err := c.ShouldBindJSON(&journal); err != nil {
		c.JSON(400, utils.NewErrorResponse(400, utils.ErrInternalServer.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(journal); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalService.CreateJournal(&ctx, &journal)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "journal created successfully", nil))
}

func (handler *AccountingHandler) UpdateJournal(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	id := c.Param("id")

	var journal accounting.AccountingJournalUpdateRequest
	if err := c.ShouldBindJSON(&journal); err != nil {
		log.Printf("%v", err)
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	if utils.IsAllFieldsNil(&journal) {
		c.JSON(400, utils.NewErrorResponse(400, "no fields to update"))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(journal); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalService.UpdateJournal(&ctx, id, &journal)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "journal updated successfully", nil))
}

func (handler *AccountingHandler) DeleteJournal(c *gin.Context) {
	id := c.Param("id")

	statusCode, err := handler.ServiceFacade.AccountingJournalService.DeleteJournal(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "journal deleted successfully", nil))
}

func (handler *AccountingHandler) JournalEntries(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingJournalEntryAllowFilterFieldsAndOps, `"accounting.journal_entry"`).
		PrepareSorts(c, accounting.AccountingJournalEntryAllowSortFields, `"limited_journal_entries"`).
		PreparePagination(c)

	journalEntries, total, statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.JournalEntries(qp)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.Header("X-Total-Count", fmt.Sprintf("%d", total))
	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"journal_entries": journalEntries,
	}))

	c.Set("total", total)
	if journalEntriesByte, err := json.Marshal(journalEntries); err == nil {
		c.Set("response", journalEntriesByte)
	}
}

func (handler *AccountingHandler) JournalEntry(c *gin.Context) {
	id := c.Param("id")

	journalEntry, statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.JournalEntry(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"journal_entry": journalEntry,
	}))
}

func (handler *AccountingHandler) CreateJournalEntry(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	var journalEntry accounting.AccountingJournalEntryCreateRequest
	if err := c.ShouldBindJSON(&journalEntry); err != nil {
		if strings.HasPrefix(err.Error(), "parsing time") {
			c.JSON(400, utils.NewErrorResponse(400, "invalid date format"))
			return
		}

		c.JSON(400, utils.NewErrorResponse(400, utils.ErrInternalServer.Error()))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(journalEntry); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.CreateJournalEntry(&ctx, &journalEntry)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "journal entry created successfully", nil))
}

func (handler *AccountingHandler) UpdateJournalEntry(c *gin.Context) {
	ctx, err := utils.Ctx(c)
	if err != nil {
		log.Printf("%v", err)
		c.JSON(500, utils.NewErrorResponse(500, utils.ErrInternalServer.Error()))
		return
	}

	id := c.Param("id")

	var journalEntry accounting.AccountingJournalEntryUpdateRequest
	if err := c.ShouldBindJSON(&journalEntry); err != nil {
		log.Printf("%v", err)
		c.JSON(400, utils.NewErrorResponse(400, utils.ErrInternalServer.Error()))
		return
	}

	if utils.IsAllFieldsNil(&journalEntry) {
		c.JSON(400, utils.NewErrorResponse(400, "no fields to update"))
		return
	}

	if validationErrors, ok := utils.ValidateStruct(journalEntry); !ok {
		c.JSON(400, utils.NewErrorResponse(400, strings.Join(validationErrors, ", ")))
		return
	}

	statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.UpdateJournalEntry(&ctx, id, &journalEntry)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "journal entry updated successfully", nil))
}

func (handler *AccountingHandler) DeleteJournalEntry(c *gin.Context) {
	id := c.Param("id")

	statusCode, err := handler.ServiceFacade.AccountingJournalEntryService.DeleteJournalEntry(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "journal entry deleted successfully", nil))
}
