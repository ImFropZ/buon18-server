package sales

import (
	"server/models"
	"strings"

	"github.com/nullism/bqb"
)

type SalesOrderItem struct {
	*models.CommonModel
	Id          int
	Name        string
	Description string
	Price       float64
	Discount    float64
}

type SalesOrderItemResponse struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	AmountTotal float64 `json:"amount_total"`
}

func SalesOrderItemToResponse(item SalesOrderItem) SalesOrderItemResponse {
	return SalesOrderItemResponse{
		Id:          item.Id,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Discount:    item.Discount,
		AmountTotal: item.Price - item.Discount,
	}
}

type SalesOrderItemCreateRequest struct {
	Name        string  `json:"name" validate:"required,max=63"`
	Description string  `json:"description" validate:"required,max=255"`
	Price       float64 `json:"price" validate:"numeric,min=0"`
	Discount    float64 `json:"discount" validate:"numeric,min=0"`
}

type SalesOrderItemUpdateRequest struct {
	Id          *int     `json:"id" validate:"required"`
	Name        *string  `json:"name" validate:"omitempty,max=63"`
	Description *string  `json:"description" validate:"omitempty,max=255"`
	Price       *float64 `json:"price" validate:"omitempty,numeric,min=0"`
	Discount    *float64 `json:"discount" validate:"omitempty,numeric,min=0"`
}

func (request SalesOrderItemUpdateRequest) MapUpdateFields(bqbQuery *bqb.Query, fieldname string, value interface{}) error {
	switch strings.ToLower(fieldname) {
	case "name":
		bqbQuery.Comma("name = ?", value)
	case "description":
		bqbQuery.Comma("description = ?", value)
	case "price":
		bqbQuery.Comma("price = ?", value)
	case "discount":
		bqbQuery.Comma("discount = ?", value)
	default:
		return models.ErrInvalidUpdateField
	}
	return nil
}
