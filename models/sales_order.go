package models

import "time"

type SalesOrder struct {
	Id           uint
	Code         string
	Note         string
	Status       string
	AcceptDate   time.Time
	DeliveryDate time.Time
	QuoteId      uint
	CId          uint
	CTime        time.Time
	MId          uint
	MTime        time.Time
}

type SalesOrderResponse struct {
	Id           uint      `json:"id"`
	Code         string    `json:"code"`
	Note         string    `json:"note"`
	Status       string    `json:"status"`
	AcceptDate   time.Time `json:"accept_date"`
	DeliveryDate time.Time `json:"delivery_date"`
	QuoteId      uint      `json:"quote_id"`
	CreatedById  uint      `json:"created_by_id"`
}

func (so *SalesOrder) PrepareForCreate(id uint) {
	so.CId = id
	so.CTime = time.Now()
	so.MId = id
	so.MTime = time.Now()
}

func (so *SalesOrder) PrepareForUpdate(id uint) {
	so.MId = id
	so.MTime = time.Now()
}

func (so *SalesOrder) ToResponse() SalesOrderResponse {
	return SalesOrderResponse{
		Id:           so.Id,
		Code:         so.Code,
		Note:         so.Note,
		Status:       so.Status,
		AcceptDate:   so.AcceptDate,
		DeliveryDate: so.DeliveryDate,
		QuoteId:      so.QuoteId,
		CreatedById:  so.CId,
	}
}

func SalesOrdersToResponse(s []SalesOrder) []SalesOrderResponse {
	res := make([]SalesOrderResponse, 0)
	for _, so := range s {
		res = append(res, so.ToResponse())
	}
	return res
}
