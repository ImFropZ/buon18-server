package accounting

import (
	"strings"

	"system.buon18.com/m/models"

	"github.com/nullism/bqb"
)

var AccountingJournalAllowFilterFieldsAndOps = []string{"code:like", "name:like", "typ:eq"}
var AccountingJournalAllowSortFields = []string{"code", "name", "typ"}

type AccountingJournal struct {
	*models.CommonModel
	Id   int
	Code string
	Name string
	Typ  string
	// -- Foreign keys
	AccountId int
}

type AccountingJournalResponse struct {
	Id      int                        `json:"id"`
	Code    string                     `json:"code"`
	Name    string                     `json:"name"`
	Typ     string                     `json:"type"`
	Account *AccountingAccountResponse `json:"account,omitempty"`
}

func AccountingJournalToResponse(
	journal AccountingJournal,
	account *AccountingAccountResponse,
) AccountingJournalResponse {
	return AccountingJournalResponse{
		Id:      journal.Id,
		Code:    journal.Code,
		Name:    journal.Name,
		Typ:     journal.Typ,
		Account: account,
	}
}

type AccountingJournalCreateRequest struct {
	Code      string `json:"code" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Typ       string `json:"type" validate:"required,accounting_journal_typ"`
	AccountId uint   `json:"account_id" validate:"required"`
}

type AccountingJournalUpdateRequest struct {
	Code      *string `json:"code"`
	Name      *string `json:"name"`
	Typ       *string `json:"type" validate:"omitempty,accounting_journal_typ"`
	AccountId *uint   `json:"account_id"`
}

func (request AccountingJournalUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "code":
		bqbQuery.Comma("code = ?", value)
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "typ":
		bqbQuery.Comma("typ = ?", value)
	case "accountid":
		bqbQuery.Comma("accounting_account_id = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}
	return nil
}
