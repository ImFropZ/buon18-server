package utils

import (
	"fmt"
	"strings"
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
	// filter = "field-operator=value"
	filterArr := strings.Split(filter, "-")
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

	if operator == "in" || operator == "nin" {
		value = fmt.Sprintf("(%s)", value)
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
