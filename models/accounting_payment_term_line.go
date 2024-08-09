package models

type AccountingPaymentTermLine struct {
	*CommonModel
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
