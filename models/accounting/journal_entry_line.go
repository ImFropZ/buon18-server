package accounting

import (
	"server/models"
	"strings"

	"github.com/nullism/bqb"
)

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

func (request AccountingJournalEntryLineUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "sequence":
		bqbQuery.Comma("sequence = ?", value)
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "amountdebit":
		bqbQuery.Comma("amount_debit = ?", value)
	case "amountcredit":
		bqbQuery.Comma("amount_credit = ?", value)
	case "accountid":
		bqbQuery.Comma("accounting_account_id = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}

	return nil
}
