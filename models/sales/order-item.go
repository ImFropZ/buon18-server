package sales

import (
	"system.buon18.com/m/models"
)

type SalesOrderItem struct {
	*models.CommonModel
	Id          *int
	Name        *string
	Description *string
	Price       *float64
	Discount    *float64
}

type SalesOrderItemResponse struct {
	Id          *int     `json:"id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Discount    *float64 `json:"discount,omitempty"`
	AmountTotal *float64 `json:"amount_total,omitempty"`
}

func SalesOrderItemToResponse(item SalesOrderItem) SalesOrderItemResponse {
	amount := 0.0
	if item.Price != nil && item.Discount != nil {
		amount = *item.Price - *item.Discount
	}

	return SalesOrderItemResponse{
		Id:          item.Id,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Discount:    item.Discount,
		AmountTotal: &amount,
	}
}

type SalesOrderItemCreateRequest struct {
	Name        string  `json:"name" validate:"required,max=63"`
	Description string  `json:"description" validate:"required,max=255"`
	Price       float64 `json:"price" validate:"numeric"`
	Discount    float64 `json:"discount" validate:"numeric,min=0"`
}

type SalesOrderItemUpdateRequest struct {
	Id          *int     `json:"id" validate:"required"`
	Name        *string  `json:"name" validate:"omitempty,max=63"`
	Description *string  `json:"description" validate:"omitempty,max=255"`
	Price       *float64 `json:"price" validate:"omitempty,numeric"`
	Discount    *float64 `json:"discount" validate:"omitempty,numeric,min=0"`
}
