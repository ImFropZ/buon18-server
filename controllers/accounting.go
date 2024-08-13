package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"server/models/accounting"
	"server/services"
	"server/utils"
	"strings"

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
