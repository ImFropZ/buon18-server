package models

import "time"

type SalesOrder struct {
	ID           uint
	Code         string
	Note         string
	Status       string
	AcceptDate   time.Time
	DeliveryDate time.Time
	QuoteID      uint
	CId          uint
	CTime        time.Time
	MId          uint
	MTime        time.Time
}

type SalesOrderResponse struct {
	ID           uint      `json:"id"`
	Code         string    `json:"code"`
	Note         string    `json:"note"`
	Status       string    `json:"status"`
	AcceptDate   time.Time `json:"accept_date"`
	DeliveryDate time.Time `json:"delivery_date"`
	QuoteID      uint      `json:"quote_id"`
	CreatedByID  uint      `json:"created_by_id"`
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
