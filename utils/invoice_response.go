package utils

import (
	"server/models"

	"github.com/gin-gonic/gin"
)

func InvoiceResponse(so models.SalesOrder, q models.Quote, c models.Client, acc models.Account, usr models.User) gin.H {
	return gin.H{
		"id":            so.Id,
		"code":          so.Code,
		"accept_date":   so.AcceptDate,
		"delivery_date": so.DeliveryDate,
		"note":          so.Note,
		"user": gin.H{
			"name":  usr.Name,
			"email": usr.Email,
			"role":  usr.Role,
		},
		"account": gin.H{
			"code":    acc.Code,
			"name":    acc.Name,
			"email":   acc.Email,
			"address": acc.Address,
			"phone":   acc.Phone,
		},
		"client": gin.H{
			"code":    c.Code,
			"name":    c.Name,
			"address": c.Address,
			"phone":   c.Phone,
		},
		"quote": gin.H{
			"code":     q.Code,
			"subtotal": q.Subtotal,
			"discount": q.Discount,
			"total":    q.Total,
			"items":    InvoiceQuoteItemsResponse(q.QuoteItems),
		},
	}
}

func InvoiceQuoteItemsResponse(qis []models.QuoteItem) []gin.H {
	res := make([]gin.H, 0)
	for _, qi := range qis {
		res = append(res, gin.H{
			"name":        qi.Name,
			"description": qi.Description,
			"quantity":    qi.Quantity,
			"unit_price":  qi.UnitPrice,
		})
	}
	return res
}
