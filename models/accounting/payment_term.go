package accounting

import (
	"server/models"
)

var AccountingPaymentTermAllowFilterFieldsAndOps = []string{"name:like", "description:like"}
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

func AccountingPaymentTermToResponse(
	term AccountingPaymentTerm,
	lines []AccountingPaymentTermLineResponse,
) AccountingPaymentTermResponse {
	return AccountingPaymentTermResponse{
		Id:          term.Id,
		Name:        term.Name,
		Description: term.Description,
		Lines:       lines,
	}
}

type AccountingPaymentTermCreateRequest struct {
	Name        string                                   `json:"name" validate:"required"`
	Description string                                   `json:"description" validate:"required"`
	Lines       []AccountingPaymentTermLineCreateRequest `json:"lines" validate:"required,gt=0,dive"`
}
