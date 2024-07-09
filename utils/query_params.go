package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationQueryParams struct {
	Offset int
	Limit  int
	Query  string
}

func (p *PaginationQueryParams) Parse(c *gin.Context) {
	if offset, err := strconv.Atoi(c.Query("offset")); err == nil {
		if offset < 0 {
			offset = 0
		}
		p.Offset = offset
	}
	if limit, err := strconv.Atoi(c.Query("limit")); err == nil {
		if limit < 1 {
			limit = 10
		}
		p.Limit = limit
	}
	if query := c.Query("q"); query != "" {
		p.Query = query
	}
}
