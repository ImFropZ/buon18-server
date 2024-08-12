package sales

import (
	"server/models"
	"server/models/setting"
	"time"
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
