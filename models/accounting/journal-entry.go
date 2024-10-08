package accounting

import (
	"strings"
	"time"

	"system.buon18.com/m/models"

	"github.com/nullism/bqb"
)

type AccountingJournalEntry struct {
	*models.CommonModel
	Id     int
	Name   string
	Date   time.Time
	Note   string
	Status string
	// -- Foreign keys
	JournalId int
}

func (AccountingJournalEntry) AllowFilterFieldsAndOps() []string {
	return []string{"status:in", "date:gte", "date:lte", "date:gt", "date:lt", "date:eq", "name:like"}
}

func (AccountingJournalEntry) AllowSorts() []string {
	return []string{"name", "date", "status"}
}

type AccountingJournalEntryResponse struct {
	Id                int                                  `json:"id"`
	Name              string                               `json:"name"`
	Date              time.Time                            `json:"date"`
	Note              string                               `json:"note"`
	Status            string                               `json:"status"`
	AmountTotalDebit  float64                              `json:"amount_total_debit"`
	AmountTotalCredit float64                              `json:"amount_total_credit"`
	Lines             []AccountingJournalEntryLineResponse `json:"lines"`
	Journal           AccountingJournalResponse            `json:"journal"`
}

func AccountingJournalEntryToResponse(
	journalEntry AccountingJournalEntry,
	lines []AccountingJournalEntryLineResponse,
	journal AccountingJournalResponse,
) AccountingJournalEntryResponse {
	response := AccountingJournalEntryResponse{
		Id:      journalEntry.Id,
		Name:    journalEntry.Name,
		Date:    journalEntry.Date,
		Note:    journalEntry.Note,
		Status:  journalEntry.Status,
		Lines:   lines,
		Journal: journal,
	}

	for _, line := range lines {
		response.AmountTotalDebit += line.AmountDebit
		response.AmountTotalCredit += line.AmountCredit
	}

	return response
}

type AccountingJournalEntryCreateRequest struct {
	Name      string                                    `json:"name" validate:"required"`
	Date      time.Time                                 `json:"date" validate:"required"`
	Note      string                                    `json:"note" validate:"required"`
	Status    string                                    `json:"status" validate:"required,accounting_journal_entry_typ"`
	JournalId int                                       `json:"journal_id" validate:"required"`
	Lines     []AccountingJournalEntryLineCreateRequest `json:"lines" validate:"required,gt=0,dive"`
}

type AccountingJournalEntryUpdateRequest struct {
	Name        *string                                    `json:"name"`
	Date        *time.Time                                 `json:"date"`
	Note        *string                                    `json:"note"`
	Status      *string                                    `json:"status" validate:"omitempty,accounting_journal_entry_typ"`
	JournalId   *int                                       `json:"journal_id"`
	AddLines    *[]AccountingJournalEntryLineCreateRequest `json:"add_lines" validate:"omitempty,gt=0,dive"`
	UpdateLines *[]AccountingJournalEntryLineUpdateRequest `json:"update_lines" validate:"omitempty,gt=0,dive"`
	DeleteLines *[]int                                     `json:"delete_lines" validate:"omitempty,gt=0,dive"`
}

func (request AccountingJournalEntryUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "date":
		bqbQuery.Comma("date = ?", value)
	case "note":
		bqbQuery.Comma("note = ?", value)
	case "status":
		bqbQuery.Comma("status = ?", value)
	case "journalid":
		bqbQuery.Comma("accounting_journal_id = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}

	return nil
}
