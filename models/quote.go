package models

import "time"

type Quote struct {
	Id         uint
	Code       string
	Date       time.Time
	ExpiryDate time.Time
	Note       string
	Subtotal   float64
	Discount   float64
	Total      float64
	ClientId   uint
	AccountId  uint
	Status     string
	CId        uint
	CTime      time.Time
	MId        uint
	MTime      time.Time

	// -- Associations
	QuoteItems []QuoteItem
}

type QuoteItem struct {
	Id          uint
	QuoteId     uint
	Name        string
	Description string
	Quantity    uint
	UnitPrice   float64
	CId         uint
	CTime       time.Time
	MId         uint
	MTime       time.Time
}

type QuoteItemCreate struct {
	Name        string
	Description *string
	Quantity    uint
	UnitPrice   float64
}

type QuoteItemUpdate struct {
	Id          uint
	Name        *string
	Description *string
	Quantity    *uint
	UnitPrice   *float64
}

type QuoteResponse struct {
	Id          uint                `json:"id"`
	Code        string              `json:"code"`
	Date        time.Time           `json:"date"`
	ExpiryDate  time.Time           `json:"expiry_date"`
	Note        string              `json:"note"`
	Subtotal    float64             `json:"subtotal"`
	Discount    float64             `json:"discount"`
	Total       float64             `json:"total"`
	Status      string              `json:"status"`
	Items       []QuoteItemResponse `json:"items"`
	ClientId    uint                `json:"client_id"`
	AccountId   uint                `json:"account_id"`
	CreatedById uint                `json:"created_by_id"`
}

type QuoteItemResponse struct {
	Id          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    uint    `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

func (q *Quote) PrepareForCreate(id uint) (err error) {
	q.CId = id
	q.CTime = time.Now()
	q.MId = id
	q.MTime = time.Now()
	return
}

func (q *Quote) PrepareForUpdate(id uint) (err error) {
	q.MId = id
	q.MTime = time.Now()
	return
}

func (quote *Quote) ToResponse() QuoteResponse {
	return QuoteResponse{
		Id:          quote.Id,
		Code:        quote.Code,
		Date:        quote.Date,
		ExpiryDate:  quote.ExpiryDate,
		Note:        quote.Note,
		Subtotal:    quote.Subtotal,
		Discount:    quote.Discount,
		Total:       quote.Total,
		Status:      quote.Status,
		Items:       QuoteItemsToResponse(quote.QuoteItems),
		ClientId:    quote.ClientId,
		AccountId:   quote.AccountId,
		CreatedById: quote.CId,
	}
}

func (quoteItem *QuoteItem) ToResponse() QuoteItemResponse {
	return QuoteItemResponse{
		Id:          quoteItem.Id,
		Name:        quoteItem.Name,
		Description: quoteItem.Description,
		Quantity:    quoteItem.Quantity,
		UnitPrice:   quoteItem.UnitPrice,
	}
}

func QuotesToResponse(quotes []Quote) []QuoteResponse {
	var response []QuoteResponse = make([]QuoteResponse, 0)

	for _, quote := range quotes {
		response = append(response, quote.ToResponse())
	}

	return response
}

func QuoteItemsToResponse(quoteItems []QuoteItem) []QuoteItemResponse {
	var response []QuoteItemResponse = make([]QuoteItemResponse, 0)

	for _, quoteItem := range quoteItems {
		response = append(response, quoteItem.ToResponse())
	}

	return response
}
