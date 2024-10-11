package accounting

import (
	"system.buon18.com/m/models"
)

type AccountingJournalEntryLine struct {
	*models.CommonModel
	Id           *int
	Sequence     *int
	Name         *string
	AmountDebit  *float64
	AmountCredit *float64
	// -- Foreign keys
	JournalEntryId *int
	AccountId      *int
}

type AccountingJournalEntryLineResponse struct {
	Id           *int                       `json:"id,omitempty"`
	Sequence     *int                       `json:"sequence,omitempty"`
	Name         *string                    `json:"name,omitempty"`
	AmountDebit  *float64                   `json:"amount_debit,omitempty"`
	AmountCredit *float64                   `json:"amount_credit,omitempty"`
	Account      *AccountingAccountResponse `json:"account,omitempty"`
}

func AccountingJournalEntryLineToResponse(line AccountingJournalEntryLine, account *AccountingAccountResponse) AccountingJournalEntryLineResponse {
	return AccountingJournalEntryLineResponse{
		Id:           line.Id,
		Sequence:     line.Sequence,
		Name:         line.Name,
		AmountDebit:  line.AmountDebit,
		AmountCredit: line.AmountCredit,
		Account:      account,
	}
}

type AccountingJournalEntryLineCreateRequest struct {
	Sequence     int     `json:"sequence" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	AmountDebit  float64 `json:"amount_debit"`
	AmountCredit float64 `json:"amount_credit"`
	AccountId    int     `json:"account_id" validate:"required"`
}

type AccountingJournalEntryLineUpdateRequest struct {
	Id           *int     `json:"id" validate:"required"`
	Sequence     *int     `json:"sequence"`
	Name         *string  `json:"name"`
	AmountDebit  *float64 `json:"amount_debit"`
	AmountCredit *float64 `json:"amount_credit"`
	AccountId    *int     `json:"account_id"`
}
