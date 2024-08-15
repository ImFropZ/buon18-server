package accounting

import (
	"server/models"
	"strings"

	"github.com/nullism/bqb"
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

type AccountingPaymentTermLineUpdateRequest struct {
	Id                 *int     `json:"id" validate:"required"`
	Sequence           *int     `json:"sequence" validate:"omitempty"`
	ValueAmountPercent *float64 `json:"value_amount_percent" validate:"omitempty"`
	NumberOfDays       *int     `json:"number_of_days" validate:"omitempty"`
}

func (request AccountingPaymentTermLineUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "sequence":
		bqbQuery.Comma("sequence = ?", value)
	case "valueamountpercent":
		bqbQuery.Comma("value_amount_percent = ?", value)
	case "numberofdays":
		bqbQuery.Comma("number_of_days = ?", value)
	}
	return nil
}
