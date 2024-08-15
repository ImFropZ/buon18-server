package sales

import (
	"server/models"
	"server/models/accounting"
	"strings"
	"time"

	"github.com/nullism/bqb"
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

func (request SalesOrderUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "commitmentdate":
		bqbQuery.Comma("commitment_date = ?", value)
	case "note":
		bqbQuery.Comma("note = ?", value)
	case "paymenttermid":
		bqbQuery.Comma("accounting_payment_term_id = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}
	return nil
}
