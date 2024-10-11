package accounting

import (
	"time"

	"system.buon18.com/m/models"
)

type AccountingJournalEntry struct {
	*models.CommonModel
	Id     *int
	Name   *string
	Date   *time.Time
	Note   *string
	Status *string
	// -- Foreign keys
	JournalId *int
}

func (AccountingJournalEntry) AllowFilterFieldsAndOps() []string {
	return []string{"status:in", "date:gte", "date:lte", "date:gt", "date:lt", "date:eq", "name:like", "name:ilike"}
}

func (AccountingJournalEntry) AllowSorts() []string {
	return []string{"name", "date", "status"}
}

type AccountingJournalEntryResponse struct {
	Id                *int                                  `json:"id,omitempty"`
	Name              *string                               `json:"name,omitempty"`
	Date              *time.Time                            `json:"date,omitempty"`
	Note              *string                               `json:"note,omitempty"`
	Status            *string                               `json:"status,omitempty"`
	AmountTotalDebit  *float64                              `json:"amount_total_debit,omitempty"`
	AmountTotalCredit *float64                              `json:"amount_total_credit,omitempty"`
	Lines             *[]AccountingJournalEntryLineResponse `json:"lines,omitempty"`
	Journal           *AccountingJournalResponse            `json:"journal,omitempty"`
}

func AccountingJournalEntryToResponse(
	journalEntry AccountingJournalEntry,
	lines *[]AccountingJournalEntryLineResponse,
	journal *AccountingJournalResponse,
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

	amountTotalCredit := 0.0
	amountTotalDebit := 0.0

	if lines != nil {
		for _, line := range *lines {
			if line.AmountDebit != nil {
				amountTotalCredit += *line.AmountDebit
			}
			if line.AmountCredit != nil {
				amountTotalDebit += *line.AmountCredit
			}
		}
	}

	response.AmountTotalCredit = &amountTotalCredit
	response.AmountTotalDebit = &amountTotalDebit

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
