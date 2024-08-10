package accounting

import "server/models"

type AccountingJournalEntryLine struct {
	*models.CommonModel
	Id           int
	Sequence     int
	Name         string
	AmountDebit  float64
	AmountCredit float64
	// -- Foreign keys
	JournalEntryId int
	AccountId      int
}

type AccountingJournalEntryLineResponse struct {
	Id           int                       `json:"id"`
	Sequence     int                       `json:"sequence"`
	Name         string                    `json:"name"`
	AmountDebit  float64                   `json:"amount_debit"`
	AmountCredit float64                   `json:"amount_credit"`
	Account      AccountingAccountResponse `json:"account"`
}

func AccountingJournalEntryLineToResponse(line AccountingJournalEntryLine, account AccountingAccountResponse) AccountingJournalEntryLineResponse {
	return AccountingJournalEntryLineResponse{
		Id:           line.Id,
		Sequence:     line.Sequence,
		Name:         line.Name,
		AmountDebit:  line.AmountDebit,
		AmountCredit: line.AmountCredit,
		Account:      account,
	}
}
