package accounting

import (
	"server/models"
	"time"
)

var AccountingJournalEntryAllowFilterFieldsAndOps = []string{"status:in", "date:gte", "date:lte", "date:gt", "date:lt", "date:eq", "name:like"}
var AccountingJournalEntryAllowSortFields = []string{"name", "date", "status"}

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
