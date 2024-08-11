package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"server/models/accounting"
	services "server/services/accounting"
	"server/utils"

	"github.com/gin-gonic/gin"
)

type AccountingHandler struct {
	DB                            *sql.DB
	AccountingAccountService      *services.AccountingAccountService
	AccountingPaymentTermService  *services.AccountingPaymentTermService
	AccountingJournalService      *services.AccountingJournalService
	AccountingJournalEntryService *services.AccountingJournalEntryService
}

func (handler *AccountingHandler) Accounts(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingAccountAllowFilterFieldsAndOps, `"accounting.account"`).
		PrepareSorts(c, accounting.AccountingAccountAllowSortFields, `"accounting.account"`).
		PreparePagination(c)

	accounts, total, statusCode, err := handler.AccountingAccountService.Accounts(qp)
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

	account, statusCode, err := handler.AccountingAccountService.Account(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"account": account,
	}))
}

func (handler *AccountingHandler) PaymentTerms(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingPaymentTermAllowFilterFieldsAndOps, `"accounting.payment_term"`).
		PrepareSorts(c, accounting.AccountingPaymentTermAllowSortFields, `"limited_payment_terms"`).
		PreparePagination(c)

	paymentTerms, total, statusCode, err := handler.AccountingPaymentTermService.PaymentTerms(qp)
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

	paymentTerm, statusCode, err := handler.AccountingPaymentTermService.PaymentTerm(id)
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

	journals, total, statusCode, err := handler.AccountingJournalService.Journals(qp)
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

	journal, statusCode, err := handler.AccountingJournalService.Journal(id)
	if err != nil {
		c.JSON(statusCode, utils.NewErrorResponse(statusCode, err.Error()))
		return
	}

	c.JSON(statusCode, utils.NewResponse(statusCode, "", gin.H{
		"journal": journal,
	}))
}

func (handler *AccountingHandler) JournalEntries(c *gin.Context) {
	qp := utils.NewQueryParams().
		PrepareFilters(c, accounting.AccountingJournalEntryAllowFilterFieldsAndOps, `"accounting.journal_entry"`).
		PrepareSorts(c, accounting.AccountingJournalEntryAllowSortFields, `"limited_journal_entries"`).
		PreparePagination(c)

	journalEntries, total, statusCode, err := handler.AccountingJournalEntryService.JournalEntries(qp)
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
