package accounting

import (
	"system.buon18.com/m/models"
)

type AccountingJournal struct {
	*models.CommonModel
	Id   *int
	Code *string
	Name *string
	Typ  *string
	// -- Foreign keys
	AccountId *int
}

func (AccountingJournal) AllowFilterFieldsAndOps() []string {
	return []string{"code:like", "code:ilike", "name:like", "name:ilike", "typ:eq"}
}

func (AccountingJournal) AllowSorts() []string {
	return []string{"code", "name", "typ"}
}

type AccountingJournalResponse struct {
	Id      *int                       `json:"id,omitempty"`
	Code    *string                    `json:"code,omitempty"`
	Name    *string                    `json:"name,omitempty"`
	Typ     *string                    `json:"type,omitempty"`
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
