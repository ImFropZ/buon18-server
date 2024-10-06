package accounting

import (
	"strings"

	"system.buon18.com/m/models"

	"github.com/nullism/bqb"
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

type AccountingPaymentTermUpdateRequest struct {
	Name          *string                                  `json:"name" validate:"omitempty"`
	Description   *string                                  `json:"description" validate:"omitempty"`
	AddLines      []AccountingPaymentTermLineCreateRequest `json:"add_lines" validate:"omitempty,dive"`
	UpdateLines   []AccountingPaymentTermLineUpdateRequest `json:"update_lines" validate:"omitempty,dive"`
	RemoveLineIds []uint                                   `json:"remove_line_ids" validate:"omitempty,dive"`
}

func (request AccountingPaymentTermUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "description":
		bqbQuery.Comma("description = ?", value)
	}
	return nil
}
