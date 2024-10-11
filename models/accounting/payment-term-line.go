package accounting

import (
	"system.buon18.com/m/models"
)

type AccountingPaymentTermLine struct {
	*models.CommonModel
	Id                 *int
	Sequence           *int
	ValueAmountPercent *float64
	NumberOfDays       *int
	// -- Foreign keys
	PaymentTermId *int
}

type AccountingPaymentTermLineResponse struct {
	Id                 *int     `json:"id,omitempty"`
	Sequence           *int     `json:"sequence,omitempty"`
	ValueAmountPercent *float64 `json:"value_amount_percent,omitempty"`
	NumberOfDays       *int     `json:"number_of_days,omitempty"`
}

func AccountingPaymentTermLineToResponse(line AccountingPaymentTermLine) AccountingPaymentTermLineResponse {
	return AccountingPaymentTermLineResponse{
		Id:                 line.Id,
		Sequence:           line.Sequence,
		ValueAmountPercent: line.ValueAmountPercent,
		NumberOfDays:       line.NumberOfDays,
	}
}

type AccountingPaymentTermLineCreateRequest struct {
	Sequence           int     `json:"sequence" validate:"required"`
	ValueAmountPercent float64 `json:"value_amount_percent" validate:"required"`
	NumberOfDays       int     `json:"number_of_days" validate:"required"`
}

type AccountingPaymentTermLineUpdateRequest struct {
	Id                 *int     `json:"id" validate:"required"`
	Sequence           *int     `json:"sequence" validate:"omitempty"`
	ValueAmountPercent *float64 `json:"value_amount_percent" validate:"omitempty"`
	NumberOfDays       *int     `json:"number_of_days" validate:"omitempty"`
}
