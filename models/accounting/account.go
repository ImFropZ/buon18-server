package accounting

import (
	"system.buon18.com/m/models"
)

type AccountingAccount struct {
	*models.CommonModel
	Id   *int
	Name *string
	Code *string
	Typ  *string
}

func (AccountingAccount) AllowFilterFieldsAndOps() []string {
	return []string{"name:like", "name:ilike", "code:like", "typ:eq"}
}

func (AccountingAccount) AllowSorts() []string {
	return []string{"name", "typ", "code"}
}

type AccountingAccountResponse struct {
	Id   *int    `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Code *string `json:"code,omitempty"`
	Typ  *string `json:"type,omitempty"`
}

func AccountingAccountToResponse(account AccountingAccount) AccountingAccountResponse {
	return AccountingAccountResponse{
		Id:   account.Id,
		Name: account.Name,
		Code: account.Code,
		Typ:  account.Typ,
	}
}

type AccountingAccountCreateRequest struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
	Typ  string `json:"type" validate:"required,accounting_account_typ"`
}

type AccountingAccountUpdateRequest struct {
	Name *string `json:"name" validate:"omitempty"`
	Code *string `json:"code" validate:"omitempty"`
	Typ  *string `json:"type" validate:"omitempty,accounting_account_typ"`
}
