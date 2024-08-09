package models

import "time"

var SalesQuotationAllowFilterFieldsAndOps = []string{"name-ilike", "status-eq", "creation_date-gte", "creation_date-lte", "validity_date-gte", "validity_date-lte"}
var SalesQuotationAllowSortFields = []string{"name", "status"}

type SalesQuotation struct {
	*CommonModel
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
	Id              int                      `json:"id"`
	Name            string                   `json:"name"`
	CreationDate    time.Time                `json:"creation_date"`
	ValidityDate    time.Time                `json:"validity_date"`
	Discount        float64                  `json:"discount"`
	AmountDelivery  float64                  `json:"amount_delivery"`
	Status          string                   `json:"status"`
	TotalAmount     float64                  `json:"total_amount"`
	Customer        SettingCustomerResponse  `json:"customer"`
	SalesOrderItems []SalesOrderItemResponse `json:"items"`
}

func SalesQuotationToResponse(quotation SalesQuotation, customer SettingCustomer, orderItems []SalesOrderItem) SalesQuotationResponse {
	subAmountTotal := 0.0
	items := make([]SalesOrderItemResponse, 0)
	for _, item := range orderItems {
		subAmountTotal += item.Price - item.Discount
		items = append(items, SalesOrderItemToResponse(item))
	}

	return SalesQuotationResponse{
		Id:              quotation.Id,
		Name:            quotation.Name,
		CreationDate:    quotation.CreationDate,
		ValidityDate:    quotation.ValidityDate,
		Discount:        quotation.Discount,
		AmountDelivery:  quotation.AmountDelivery,
		Status:          quotation.Status,
		Customer:        SettingCustomerToResponse(customer),
		TotalAmount:     subAmountTotal + quotation.AmountDelivery - quotation.Discount,
		SalesOrderItems: items,
	}
}
