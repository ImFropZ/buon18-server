package accounting

import "server/models"

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
	Id      int                       `json:"id"`
	Code    string                    `json:"code"`
	Name    string                    `json:"name"`
	Typ     string                    `json:"type"`
	Account AccountingAccountResponse `json:"account"`
}

func AccountingJournalToResponse(
	journal AccountingJournal,
	account AccountingAccountResponse,
) AccountingJournalResponse {
	return AccountingJournalResponse{
		Id:      journal.Id,
		Code:    journal.Code,
		Name:    journal.Name,
		Typ:     journal.Typ,
		Account: account,
	}
}
