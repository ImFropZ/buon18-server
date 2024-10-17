package sales

import (
	"time"

	"system.buon18.com/m/models"
	"system.buon18.com/m/models/accounting"
)

type SalesOrder struct {
	*models.CommonModel
	Id             *int
	Name           *string
	CommitmentDate *time.Time
	Note           *string
	// -- Foreign keys
	SalesQuotationId        *int
	AccountingPaymentTermId *int
}

func (SalesOrder) AllowFilterFieldsAndOps() []string {
	return []string{"id:eq", "name:like", "commitment-date:eq", "commitment-date:gt", "commitment-date:gte", "commitment-date:lt", "commitment-date:lte", "sales-quotation-id:eq", "accounting-payment-term-id:eq"}
}

func (SalesOrder) AllowSorts() []string {
	return []string{"name", "commitment-date"}
}

type SalesOrderResponse struct {
	Id             *int                                      `json:"id,omitempty"`
	Name           *string                                   `json:"name,omitempty"`
	CommitmentDate *time.Time                                `json:"commitment_date,omitempty"`
	Note           *string                                   `json:"note,omitempty"`
	Quotation      *SalesQuotationResponse                   `json:"quotation,omitempty"`
	PaymentTerm    *accounting.AccountingPaymentTermResponse `json:"payment_term,omitempty"`
}

func SalesOrderToResponse(
	salesOrder SalesOrder,
	quotation *SalesQuotationResponse,
	paymentTerm *accounting.AccountingPaymentTermResponse,
) SalesOrderResponse {
	return SalesOrderResponse{
		Id:             salesOrder.Id,
		Name:           salesOrder.Name,
		CommitmentDate: salesOrder.CommitmentDate,
		Note:           salesOrder.Note,
		Quotation:      quotation,
		PaymentTerm:    paymentTerm,
	}
}

type SalesOrderCreateRequest struct {
	Name           string    `json:"name" validate:"required"`
	CommitmentDate time.Time `json:"commitment_date" validate:"required"`
	Note           string    `json:"note"`
	QuotationId    uint      `json:"quotation_id" validate:"required"`
	PaymentTermId  uint      `json:"payment_term_id" validate:"required"`
}

type SalesOrderUpdateRequest struct {
	Name           *string    `json:"name" validate:"omitempty"`
	CommitmentDate *time.Time `json:"commitment_date" validate:"omitempty"`
	Note           *string    `json:"note" validate:"omitempty"`
	PaymentTermId  *uint      `json:"payment_term_id" validate:"omitempty"`
}
