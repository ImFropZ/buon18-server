package accounting

import (
	"server/models"
)

var AccountingPaymentTermAllowFilterFieldsAndOps = []string{"name-like", "description-like"}
var AccountingPaymentTermAllowSortFields = []string{"name"}

type AccountingPaymentTerm struct {
	*models.CommonModel
	Id          int
	Name        string
	Description string
}

type AccountingPaymentTermResponse struct {
	Id          int                                 `json:"id"`
	Name        string                              `json:"name"`
	Description string                              `json:"description"`
	Lines       []AccountingPaymentTermLineResponse `json:"lines"`
}

func AccountingPaymentTermToResponse(term AccountingPaymentTerm, paymentTermLines []AccountingPaymentTermLine) AccountingPaymentTermResponse {
	lines := make([]AccountingPaymentTermLineResponse, 0)
	for _, line := range paymentTermLines {
		lines = append(lines, AccountingPaymentTermLineToResponse(line))
	}
	return AccountingPaymentTermResponse{
		Id:          term.Id,
		Name:        term.Name,
		Description: term.Description,
		Lines:       lines,
	}
}
