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
	Client     Client
	Account    Account
	QuoteItems []QuoteItem
	CreatedBy  User
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

type QuoteResponse struct {
	Id         uint                `json:"id"`
	Code       string              `json:"code"`
	Date       time.Time           `json:"date"`
	ExpiryDate time.Time           `json:"expiry_date"`
	Note       string              `json:"note"`
	Subtotal   float64             `json:"subtotal"`
	Discount   float64             `json:"discount"`
	Total      float64             `json:"total"`
	Status     string              `json:"status"`
	Items      []QuoteItemResponse `json:"items"`
	Client     ClientResponse      `json:"client"`
	Account    AccountResponse     `json:"account"`
	CreatedBy  UserResponse        `json:"created_by"`
}

type QuoteItemResponse struct {
	Id          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    uint    `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

func (quote *Quote) ToResponse() QuoteResponse {
	return QuoteResponse{
		Id:         quote.Id,
		Code:       quote.Code,
		Date:       quote.Date,
		ExpiryDate: quote.ExpiryDate,
		Note:       quote.Note,
		Subtotal:   quote.Subtotal,
		Discount:   quote.Discount,
		Total:      quote.Total,
		Status:     quote.Status,
		Items:      QuoteItemsToResponse(quote.QuoteItems),
		Client:     quote.Client.ToResponse(),
		Account:    quote.Account.ToResponse(),
		CreatedBy:  quote.CreatedBy.ToResponse(),
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
