package sales

import (
	"strings"
	"time"

	"system.buon18.com/m/models"
	"system.buon18.com/m/models/setting"

	"github.com/nullism/bqb"
)

var SalesQuotationAllowFilterFieldsAndOps = []string{"name:like", "status:eq", "creation_date:gte", "creation_date:lte", "validity_date:gte", "validity_date:lte"}
var SalesQuotationAllowSortFields = []string{"name", "status"}

type SalesQuotation struct {
	*models.CommonModel
	Id             int
	Name           string
	CreationDate   time.Time
	ValidityDate   time.Time
	Discount       float64
	AmountDelivery float64
	Status         string
	// -- Foreign keys
	CustomerId int
}

type SalesQuotationResponse struct {
	Id              int                             `json:"id"`
	Name            string                          `json:"name"`
	CreationDate    time.Time                       `json:"creation_date"`
	ValidityDate    time.Time                       `json:"validity_date"`
	Discount        float64                         `json:"discount"`
	AmountDelivery  float64                         `json:"amount_delivery"`
	Status          string                          `json:"status"`
	TotalAmount     float64                         `json:"total_amount"`
	Customer        setting.SettingCustomerResponse `json:"customer"`
	SalesOrderItems []SalesOrderItemResponse        `json:"items"`
}

func SalesQuotationToResponse(
	quotation SalesQuotation,
	customer setting.SettingCustomerResponse,
	orderItems []SalesOrderItemResponse,
) SalesQuotationResponse {
	subAmountTotal := 0.0
	for _, item := range orderItems {
		subAmountTotal += item.AmountTotal
	}

	return SalesQuotationResponse{
		Id:              quotation.Id,
		Name:            quotation.Name,
		CreationDate:    quotation.CreationDate,
		ValidityDate:    quotation.ValidityDate,
		Discount:        quotation.Discount,
		AmountDelivery:  quotation.AmountDelivery,
		Status:          quotation.Status,
		Customer:        customer,
		TotalAmount:     subAmountTotal + quotation.AmountDelivery - quotation.Discount,
		SalesOrderItems: orderItems,
	}
}

type SalesQuotationCreateRequest struct {
	Name            string                        `json:"name" validate:"required"`
	CreationDate    time.Time                     `json:"creation_date" validate:"required"`
	ValidityDate    time.Time                     `json:"validity_date" validate:"required"`
	Discount        float64                       `json:"discount" validate:"numeric,min=0"`
	AmountDelivery  float64                       `json:"amount_delivery" validate:"numeric,min=0"`
	Status          string                        `json:"status" validate:"required,sales_quotation_status"`
	CustomerId      uint                          `json:"customer_id" validate:"required"`
	SalesOrderItems []SalesOrderItemCreateRequest `json:"items" validate:"required,gt=0,dive"`
}

type SalesQuotationUpdateRequest struct {
	Name                    *string                        `json:"name" validate:"omitempty"`
	CreationDate            *time.Time                     `json:"creation_date" validate:"omitempty"`
	ValidityDate            *time.Time                     `json:"validity_date" validate:"omitempty"`
	Discount                *float64                       `json:"discount" validate:"omitempty,numeric,min=0"`
	AmountDelivery          *float64                       `json:"amount_delivery" validate:"omitempty,numeric,min=0"`
	Status                  *string                        `json:"status" validate:"omitempty,sales_quotation_status"`
	CustomerId              *uint                          `json:"customer_id" validate:"omitempty"`
	AddSalesOrderItems      *[]SalesOrderItemCreateRequest `json:"add_items" validate:"omitempty,gt=0,dive"`
	UpdateSalesOrderItems   *[]SalesOrderItemUpdateRequest `json:"update_items" validate:"omitempty,gt=0,dive"`
	DeleteSalesOrderItemIds *[]uint                        `json:"delete_item_ids" validate:"omitempty,gt=0,dive"`
}

func (request SalesQuotationUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "creation_date":
		bqbQuery.Comma("creation_date = ?", value)
	case "validity_date":
		bqbQuery.Comma("validity_date = ?", value)
	case "discount":
		bqbQuery.Comma("discount = ?", value)
	case "amount_delivery":
		bqbQuery.Comma("amount_delivery = ?", value)
	case "status":
		bqbQuery.Comma("status = ?", value)
	case "customer_id":
		bqbQuery.Comma("customer_id = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}
	return nil
}
