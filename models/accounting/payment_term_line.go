package accounting

import (
	"server/models"
)

type AccountingPaymentTermLine struct {
	*models.CommonModel
	Id                 int
	Sequence           int
	ValueAmountPercent float64
	NumberOfDays       int
	// -- Foreign keys
	PaymentTermId int
}

type AccountingPaymentTermLineResponse struct {
	Id                 int     `json:"id"`
	Sequence           int     `json:"sequence"`
	ValueAmountPercent float64 `json:"value_amount_percent"`
	NumberOfDays       int     `json:"number_of_days"`
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
