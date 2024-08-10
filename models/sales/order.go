package sales

import (
	"server/models"
	"server/models/accounting"
	"server/models/setting"
	"time"
)

var SalesOrderAllowFilterFieldsAndOps = []string{}
var SalesOrderAllowSortFields = []string{}

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
	quotation SalesQuotation,
	customer setting.SettingCustomer,
	orderItems []SalesOrderItem,
	paymentTerm accounting.AccountingPaymentTerm,
	paymentTermLines []accounting.AccountingPaymentTermLine,
) SalesOrderResponse {
	return SalesOrderResponse{
		Id:             salesOrder.Id,
		Name:           salesOrder.Name,
		CommitmentDate: salesOrder.CommitmentDate,
		Note:           salesOrder.Note,
		Quotation:      SalesQuotationToResponse(quotation, customer, orderItems),
		PaymentTerm:    accounting.AccountingPaymentTermToResponse(paymentTerm, paymentTermLines),
	}
}
