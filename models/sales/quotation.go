package sales

import (
	"time"

	"system.buon18.com/m/models"
	"system.buon18.com/m/models/setting"
)

type SalesQuotation struct {
	*models.CommonModel
	Id             *int
	Name           *string
	CreationDate   *time.Time
	ValidityDate   *time.Time
	Discount       *float64
	AmountDelivery *float64
	Status         *string
	// -- Foreign keys
	CustomerId *int
}

func (SalesQuotation) AllowFilterFieldsAndOps() []string {
	return []string{"id:eq", "name:like", "name:ilike", "status:eq", "creation-date:gte", "creation-date:lte", "validity-date:gte", "validity-date:lte"}
}

func (SalesQuotation) AllowSorts() []string {
	return []string{"name", "status", "creation-date", "validity-date"}
}

type SalesQuotationResponse struct {
	Id              *int                             `json:"id,omitempty"`
	Name            *string                          `json:"name,omitempty"`
	CreationDate    *time.Time                       `json:"creation_date,omitempty"`
	ValidityDate    *time.Time                       `json:"validity_date,omitempty"`
	Discount        *float64                         `json:"discount,omitempty"`
	AmountDelivery  *float64                         `json:"amount_delivery,omitempty"`
	Status          *string                          `json:"status,omitempty"`
	TotalAmount     *float64                         `json:"total_amount,omitempty"`
	Customer        *setting.SettingCustomerResponse `json:"customer,omitempty"`
	SalesOrderItems *[]SalesOrderItemResponse        `json:"items,omitempty"`
}

func SalesQuotationToResponse(
	quotation SalesQuotation,
	customer *setting.SettingCustomerResponse,
	orderItems *[]SalesOrderItemResponse,
) SalesQuotationResponse {
	subAmountTotal := 0.0
	if orderItems != nil {
		for _, item := range *orderItems {
			if item.AmountTotal != nil {
				subAmountTotal += *item.AmountTotal
			}
		}
	}

	discount := 0.0
	if quotation.Discount != nil {
		discount = *quotation.Discount
	}

	amountDelivery := 0.0
	if quotation.AmountDelivery != nil {
		amountDelivery = *quotation.AmountDelivery
	}

	total := subAmountTotal + amountDelivery - discount

	return SalesQuotationResponse{
		Id:              quotation.Id,
		Name:            quotation.Name,
		CreationDate:    quotation.CreationDate,
		ValidityDate:    quotation.ValidityDate,
		Discount:        quotation.Discount,
		AmountDelivery:  quotation.AmountDelivery,
		Status:          quotation.Status,
		Customer:        customer,
		TotalAmount:     &total,
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
