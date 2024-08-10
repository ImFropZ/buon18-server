package utils

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nullism/bqb"
	"github.com/stretchr/testify/assert"
)

func TestQueryParams(t *testing.T) {
	t.Run("Create", func(t *testing.T) {
		qp := NewQueryParams()
		if len(qp.Fitlers) != 0 {
			t.Errorf("Expected 0, got %d", len(qp.Fitlers))
		}
		if qp.Pagination.Offset != 0 {
			t.Errorf("Expected 0, got %d", qp.Pagination.Offset)
		}
		if qp.Pagination.Limit != 10 {
			t.Errorf("Expected 10, got %d", qp.Pagination.Limit)
		}
		if len(qp.OrderBy) != 0 {
			t.Errorf("Expected 0, got %d", len(qp.OrderBy))
		}
	})

	t.Run("AddFilter", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddFilter("field:eq=value")
		if len(qp.Fitlers) != 1 {
			t.Errorf("Expected 1, got %d", len(qp.Fitlers))
		}
		if qp.Fitlers[0].Field != "field" {
			t.Errorf("Expected field, got %s", qp.Fitlers[0].Field)
		}
		if qp.Fitlers[0].Operator != "eq" {
			t.Errorf("Expected eq, got %s", qp.Fitlers[0].Operator)
		}
		if qp.Fitlers[0].Value != "value" {
			t.Errorf("Expected value, got %s", qp.Fitlers[0].Value)
		}
	})

	t.Run("AddFilterWithLikeOperator", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddFilter("field:like=value")
		if len(qp.Fitlers) != 1 {
			t.Errorf("Expected 1, got %d", len(qp.Fitlers))
		}
		if qp.Fitlers[0].Field != "field" {
			t.Errorf("Expected field, got %s", qp.Fitlers[0].Field)
		}
		if qp.Fitlers[0].Operator != "like" {
			t.Errorf("Expected in, got %s", qp.Fitlers[0].Operator)
		}
		if qp.Fitlers[0].Value != "%value%" {
			t.Errorf(`Expected %%value%%, got %s`, qp.Fitlers[0].Value)
		}
	})

	t.Run("AddFilterWithInvalidOperator", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddFilter("field:invalid=value")
		if len(qp.Fitlers) != 0 {
			t.Errorf("Expected 0, got %d", len(qp.Fitlers))
		}
	})

	t.Run("Pagination", func(t *testing.T) {
		qp := NewQueryParams()
		qp.Pagination.Offset = 5
		qp.Pagination.Limit = 20
		if qp.Pagination.Offset != 5 {
			t.Errorf("Expected 5, got %d", qp.Pagination.Offset)
		}
		if qp.Pagination.Limit != 20 {
			t.Errorf("Expected 20, got %d", qp.Pagination.Limit)
		}
	})

	t.Run("PaginationWithInvalidOffsetOrLimit", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddOffset(-5)
		if qp.Pagination.Offset != 0 {
			t.Errorf("Expected 0, got %d", qp.Pagination.Offset)
		}
		if qp.Pagination.Limit != 10 {
			t.Errorf("Expected 10, got %d", qp.Pagination.Limit)
		}

		qp = NewQueryParams()
		qp.AddLimit(-20)
		if qp.Pagination.Offset != 0 {
			t.Errorf("Expected 0, got %d", qp.Pagination.Offset)
		}
		if qp.Pagination.Limit != 10 {
			t.Errorf("Expected 10, got %d", qp.Pagination.Limit)
		}
	})

	t.Run("AddOrderBy", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddOrderBy("field1 ASC")
		qp.AddOrderBy("field2 DESC")
		if len(qp.OrderBy) != 2 {
			t.Errorf("Expected 2, got %d", len(qp.OrderBy))
		}
		if qp.OrderBy[0] != "field1 ASC" {
			t.Errorf("Expected field1, got %s", qp.OrderBy[0])
		}
		if qp.OrderBy[1] != "field2 DESC" {
			t.Errorf("Expected field2, got %s", qp.OrderBy[1])
		}
	})

	t.Run("FilterIntoBqb", func(t *testing.T) {
		bqbQuery := bqb.New("SELECT * FROM table")
		qp := NewQueryParams()

		qp.AddFilter("field1:eq=value1")
		qp.AddFilter("field2:in=value2,value3")

		qp.FilterIntoBqb(bqbQuery)

		expectedQuery := `SELECT * FROM table WHERE field1 = $1 AND field2 IN ( $2 , $3 )`
		expectedParams := []interface{}{"value1", "value2", "value3"}
		query, params, err := bqbQuery.ToPgsql()
		assert.NoError(t, err)

		if query != expectedQuery {
			t.Errorf("Expected %s, got %s", expectedQuery, query)
		}

		if params[0] != expectedParams[0] || params[1] != expectedParams[1] || params[2] != expectedParams[2] {
			t.Errorf("Expected %v, got %v", expectedParams, params)
		}
	})

	t.Run("PaginationIntoBqb", func(t *testing.T) {
		bqbQuery := bqb.New("SELECT * FROM table")
		qp := NewQueryParams()

		qp.Pagination.Offset = 5
		qp.Pagination.Limit = 10

		qp.PaginationIntoBqb(bqbQuery)

		expectedQuery := `SELECT * FROM table OFFSET $1 LIMIT $2`
		expectedParams := []interface{}{5, 10}
		query, params, err := bqbQuery.ToPgsql()
		assert.NoError(t, err)

		if query != expectedQuery {
			t.Errorf("Expected %s, got %s", expectedQuery, query)
		}

		if params[0] != expectedParams[0] || params[1] != expectedParams[1] {
			t.Errorf("Expected %v, got %v", expectedParams, params)
		}
	})

	t.Run("OrderByIntoBqb", func(t *testing.T) {
		bqbQuery := bqb.New("SELECT * FROM table")
		qp := NewQueryParams()

		qp.AddOrderBy("field1 ASC")
		qp.AddOrderBy("field2 DESC")

		qp.OrderByIntoBqb(bqbQuery, "field3 ASC")

		expectedQuery := `SELECT * FROM table ORDER BY field1 ASC , field2 DESC , field3 ASC`
		query, _, err := bqbQuery.ToPgsql()
		assert.NoError(t, err)

		if query != expectedQuery {
			t.Errorf("Expected %s, got %s", expectedQuery, query)
		}
	})

	t.Run("PrepareFitler", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/?field1:eq=foo&field2:in=bar,baz", nil)

		qp := NewQueryParams().
			PrepareFilters(c, []string{"field1:eq", "field2:in"}, `"foo.bar"`)

		if len(qp.Fitlers) != 2 {
			t.Errorf("Expected 2, got %d", len(qp.Fitlers))
		}

		if qp.Fitlers[0].Field != `"foo.bar".field1` {
			t.Errorf(`Expected "foo.bar".field1, got %s`, qp.Fitlers[0].Field)
		}
		if qp.Fitlers[0].Operator != "eq" {
			t.Errorf("Expected eq, got %s", qp.Fitlers[0].Operator)
		}
		if qp.Fitlers[0].Value != "foo" {
			t.Errorf("Expected foo, got %s", qp.Fitlers[0].Value)
		}

		if qp.Fitlers[1].Field != `"foo.bar".field2` {
			t.Errorf(`Expected "foo.bar".field2, got %s`, qp.Fitlers[1].Field)
		}
		if qp.Fitlers[1].Operator != "in" {
			t.Errorf("Expected in, got %s", qp.Fitlers[1].Operator)
		}
		if qp.Fitlers[1].Value != "bar,baz" {
			t.Errorf("Expected bar,baz got %s", qp.Fitlers[1].Value)
		}
	})

	t.Run("PrepareSort", func(t *testing.T) {
		gin.SetMode(gin.TestMode)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("GET", "/?sort:email=asc&sort:name=desc", nil)

		qp := NewQueryParams().
			PrepareSorts(c, []string{"email", "name"}, `"foo.bar"`)

		if len(qp.OrderBy) != 2 {
			t.Errorf("Expected 2, got %d", len(qp.OrderBy))
		}

		if qp.OrderBy[0] != `LOWER("foo.bar".email) ASC` {
			t.Errorf(`Expected LOWER("foo.bar".email) ASC, got %s`, qp.OrderBy[0])
		}

		if qp.OrderBy[1] != `LOWER("foo.bar".name) DESC` {
			t.Errorf(`Expected LOWER("foo.bar".name) DESC, got %s`, qp.OrderBy[1])
		}
	})

}
