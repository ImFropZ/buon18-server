package accounting

import (
	"system.buon18.com/m/models"
)

type AccountingPaymentTerm struct {
	*models.CommonModel
	Id          *int
	Name        *string
	Description *string
}

func (AccountingPaymentTerm) AllowFilterFieldsAndOps() []string {
	return []string{"name:like", "name:ilike", "description:like", "description:ilike"}
}

func (AccountingPaymentTerm) AllowSorts() []string {
	return []string{"name"}
}

type AccountingPaymentTermResponse struct {
	Id          *int                                 `json:"id,omitempty"`
	Name        *string                              `json:"name,omitempty"`
	Description *string                              `json:"description,omitempty"`
	Lines       *[]AccountingPaymentTermLineResponse `json:"lines,omitempty"`
}

func AccountingPaymentTermToResponse(
	term AccountingPaymentTerm,
	lines *[]AccountingPaymentTermLineResponse,
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
