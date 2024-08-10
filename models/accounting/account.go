package accounting

import "server/models"

var AccountingAccountAllowFilterFieldsAndOps = []string{"name:like", "code:like", "typ:eq"}
var AccountingAccountAllowSortFields = []string{"name", "typ"}

type AccountingAccount struct {
	*models.CommonModel
	Id   int
	Name string
	Code string
	Typ  string
}

type AccountingAccountResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Typ  string `json:"type"`
}

func AccountingAccountToResponse(account AccountingAccount) AccountingAccountResponse {
	return AccountingAccountResponse{
		Id:   account.Id,
		Name: account.Name,
		Code: account.Code,
		Typ:  account.Typ,
	}
}
