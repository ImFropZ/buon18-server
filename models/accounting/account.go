package accounting

import (
	"strings"

	"system.buon18.com/m/models"

	"github.com/nullism/bqb"
)

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

func (request AccountingAccountUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "code":
		bqbQuery.Comma("code = ?", value)
	case "type":
		bqbQuery.Comma("typ = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}
	return nil
}
