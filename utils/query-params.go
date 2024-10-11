package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/nullism/bqb"
	"system.buon18.com/m/models"
)

var ALLOWED_FILTER_OPERATORS = []string{"eq", "ne", "gt", "lt", "gte", "lte", "like", "ilike", "in", "nin"}

var MAPPED_FILTER_OPERATORS_TO_SQL = map[string]string{
	"eq":    "=",
	"ne":    "!=",
	"gt":    ">",
	"lt":    "<",
	"gte":   ">=",
	"lte":   "<=",
	"like":  "LIKE",
	"ilike": "ILIKE",
	"in":    "IN",
	"nin":   "NOT IN",
}

type FilterValue struct {
	Field    string
	Operator string
	Value    string
}

type QueryParams struct {
	Fitlers []FilterValue
	OrderBy []string
	Offset  int
	Limit   int
}

func NewQueryParams() *QueryParams {
	return &QueryParams{
		Fitlers: make([]FilterValue, 0),
		OrderBy: make([]string, 0),
		Offset:  0,
		Limit:   10,
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

	if operator == "like" || operator == "ilike" {
		value = "%" + value + "%"
	}

	qp.Fitlers = append(qp.Fitlers, FilterValue{field, operator, value})
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

	// NOTE: replace `-` with `_` in field name to match with database column name and clean query
	field = strings.ReplaceAll(field, `-`, "_")
	qp.OrderBy = append(qp.OrderBy, fmt.Sprintf("%s %s", field, strings.ToUpper(sort)))
	return qp
}

func (qp *QueryParams) FilterIntoBqb(bqbQuery *bqb.Query) {
	if len(qp.Fitlers) > 0 {
		bqbQuery.Space("WHERE")
		for index, filter := range qp.Fitlers {
			// NOTE: replace `-` with `_` in field name to match with database column name and clean query
			filter.Field = strings.ReplaceAll(filter.Field, `-`, "_")
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
		bqbQuery.Space(`ORDER BY ` + defaultOrderBy)
	}
}

func (qp *QueryParams) PaginationIntoBqb(bqbQuery *bqb.Query) {
	bqbQuery.Space(`OFFSET ? LIMIT ?`, qp.Offset, qp.Limit)
}

func (qp *QueryParams) PrepareFilters(model models.IFilter, r *http.Request, prefix string) *QueryParams {
	c := r.URL.Query()
	for _, filter := range model.AllowFilterFieldsAndOps() {
		if q := c.Get(filter); q != "" {
			qp.AddFilter(fmt.Sprintf(`%s.%s=%s`, prefix, filter, q))
		}
	}
	return qp
}

func (qp *QueryParams) PrepareSorts(model models.ISort, r *http.Request, prefix string) *QueryParams {
	c := r.URL.Query()
	for _, sort := range model.AllowSorts() {
		if q := c.Get(fmt.Sprintf("sort:%s", sort)); q != "" {
			qp.AddOrderBy(fmt.Sprintf(`LOWER("%s".%s) %s`, prefix, sort, q))
		}
	}
	return qp
}

func (qp *QueryParams) PrepareLimitAndOffset(r *http.Request) *QueryParams {
	c := r.URL.Query()
	if q := c.Get("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil {
			qp.Limit = n
		}
	}
	if q := c.Get("offset"); q != "" {
		if n, err := strconv.Atoi(q); err == nil {
			qp.Offset = n
		}
	}
	return qp
}
