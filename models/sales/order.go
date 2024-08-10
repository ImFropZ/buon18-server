package sales

import (
	"server/models"
	"server/models/accounting"
	"time"
)

var SalesOrderAllowFilterFieldsAndOps = []string{"name:like", "commitment_date:eq", "commitment_date:gt", "commitment_date:gte", "commitment_date:lt", "commitment_date:lte", "sales_quotation_id:eq", "accounting_payment_term_id:eq"}
var SalesOrderAllowSortFields = []string{"name", "commitment_date"}

type SalesOrder struct {
	*models.CommonModel
	Id             int
	Name           string
	CommitmentDate time.Time
	Note           string
	// -- Foreign keys
	SalesQuotationId        int
	AccountingPaymentTermId int
}

type SalesOrderResponse struct {
	Id             int                                      `json:"id"`
	Name           string                                   `json:"name"`
	CommitmentDate time.Time                                `json:"commitment_date"`
	Note           string                                   `json:"note"`
	Quotation      SalesQuotationResponse                   `json:"quotation"`
	PaymentTerm    accounting.AccountingPaymentTermResponse `json:"payment_term"`
}

func SalesOrderToResponse(
	salesOrder SalesOrder,
	quotation SalesQuotationResponse,
	paymentTerm accounting.AccountingPaymentTermResponse,
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
