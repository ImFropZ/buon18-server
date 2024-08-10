package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
)

var ALLOWED_FILTER_OPERATORS = []string{"eq", "ne", "gt", "lt", "gte", "lte", "like", "in", "nin"}

var MAPPED_FILTER_OPERATORS_TO_SQL = map[string]string{
	"eq":   "=",
	"ne":   "!=",
	"gt":   ">",
	"lt":   "<",
	"gte":  ">=",
	"lte":  "<=",
	"like": "LIKE",
	"in":   "IN",
	"nin":  "NOT IN",
}

type FilterValue struct {
	Field    string
	Operator string
	Value    string
}

type PaginationValue struct {
	Offset int
	Limit  int
}

type QueryParams struct {
	Fitlers    []FilterValue
	Pagination PaginationValue
	OrderBy    []string
}

func NewQueryParams() *QueryParams {
	return &QueryParams{
		Fitlers:    []FilterValue{},
		Pagination: PaginationValue{0, 10},
		OrderBy:    []string{},
	}
}

func (qp *QueryParams) AddFilter(filter string) *QueryParams {
	// filter = "field:operator=value"
	filterArr := strings.Split(filter, ":")
	field := filterArr[0]
	if len(filterArr) != 2 {
		return qp
	}
	filterArr = strings.Split(filterArr[1], "=")
	if len(filterArr) != 2 {
		return qp
	}
	operator := filterArr[0]
	value := filterArr[1]

	// -- validate operator
	if !ContainsString(ALLOWED_FILTER_OPERATORS, operator) {
		return qp
	}

	if operator == "like" {
		value = "%" + value + "%"
	}

	qp.Fitlers = append(qp.Fitlers, FilterValue{field, operator, value})
	return qp
}

func (qp *QueryParams) AddOffset(offset int) *QueryParams {
	if offset < 0 {
		return qp
	}

	qp.Pagination = PaginationValue{offset, qp.Pagination.Limit}
	return qp
}

func (qp *QueryParams) AddLimit(limit int) *QueryParams {
	if limit < 0 {
		return qp
	}

	qp.Pagination = PaginationValue{qp.Pagination.Offset, limit}
	return qp
}

func (qp *QueryParams) AddOrderBy(orderBy string) *QueryParams {
	orderByArr := strings.Split(orderBy, " ")
	if len(orderByArr) != 2 {
		return qp
	}

	field := orderByArr[0]
	sort := orderByArr[1]

	if !(strings.EqualFold(sort, "asc") || strings.EqualFold(sort, "desc")) {
		return qp
	}

	qp.OrderBy = append(qp.OrderBy, fmt.Sprintf("%s %s", field, strings.ToUpper(sort)))
	return qp
}

func (qp *QueryParams) FilterIntoBqb(bqbQuery *bqb.Query) {
	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			if filter.Operator == "in" || filter.Operator == "nin" {
				values := strings.Split(filter.Value, ",")
				bqbQuery.Space(fmt.Sprintf("%s %s (", filter.Field, MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]))
				for i, value := range values {
					bqbQuery.Space("?", value)
					if i < len(values)-1 {
						bqbQuery.Space(",")
					}
				}
				bqbQuery.Space(")")
			} else {
				bqbQuery.Space(fmt.Sprintf("%s %s ?", filter.Field, MAPPED_FILTER_OPERATORS_TO_SQL[filter.Operator]), filter.Value)
			}
			if index < len(qp.Fitlers)-1 {
				bqbQuery.Space("AND")
			}
		}
	}
}

func (qp *QueryParams) OrderByIntoBqb(bqbQuery *bqb.Query, defaultOrderBy string) {
	if len(qp.OrderBy) > 0 {
		bqbQuery.Space("ORDER BY")
		for index, sort := range qp.OrderBy {
			bqbQuery.Space(sort)
			if index < len(qp.OrderBy)-1 {
				bqbQuery.Space(",")
			}
		}
		if defaultOrderBy != "" {
			bqbQuery.Space(fmt.Sprintf(", %s", defaultOrderBy))
		}
	} else {
		bqbQuery.Space(`ORDER BY` + defaultOrderBy)
	}
}

func (qp *QueryParams) PaginationIntoBqb(bqbQuery *bqb.Query) {
	bqbQuery.Space(`OFFSET ? LIMIT ?`, qp.Pagination.Offset, qp.Pagination.Limit)
}

func (qp *QueryParams) PrepareFilters(c *gin.Context, allowFilterFieldsAndOps []string, prefix string) *QueryParams {
	for _, filter := range allowFilterFieldsAndOps {
		if validFilter, ok := c.GetQuery(filter); ok {
			qp.AddFilter(fmt.Sprintf(`%s.%s=%s`, prefix, filter, validFilter))
		}
	}
	return qp
}

func (qp *QueryParams) PrepareSorts(c *gin.Context, allowSortFields []string, prefix string) *QueryParams {
	for _, sort := range allowSortFields {
		if validSort, ok := c.GetQuery(fmt.Sprintf("sort:%s", sort)); ok {
			qp.AddOrderBy(fmt.Sprintf(`LOWER(%s.%s) %s`, prefix, sort, validSort))
		}
	}
	return qp
}

func (qp *QueryParams) PreparePagination(c *gin.Context) *QueryParams {
	for _, pagination := range []string{"offset", "limit"} {
		if validPagination, ok := c.GetQuery(pagination); ok {
			if pagination == "offset" {
				qp.AddOffset(StrToInt(validPagination, 0))
			} else {
				qp.AddLimit(StrToInt(validPagination, 10))
			}
		}
	}
	return qp
}
