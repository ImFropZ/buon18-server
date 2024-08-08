package utils

import "testing"

func TestFilterFeatures(t *testing.T) {
	t.Run("Test NewQueryParams", func(t *testing.T) {
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

	t.Run("Test AddFilter", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddFilter("field-eq=value")
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

	t.Run("Test AddFilter with `in` operator", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddFilter("field-in=value1,value2")
		if len(qp.Fitlers) != 1 {
			t.Errorf("Expected 1, got %d", len(qp.Fitlers))
		}
		if qp.Fitlers[0].Field != "field" {
			t.Errorf("Expected field, got %s", qp.Fitlers[0].Field)
		}
		if qp.Fitlers[0].Operator != "in" {
			t.Errorf("Expected in, got %s", qp.Fitlers[0].Operator)
		}
		if qp.Fitlers[0].Value != "(value1,value2)" {
			t.Errorf("Expected (value1,value2), got %s", qp.Fitlers[0].Value)
		}
	})

	t.Run("Test AddFilter with `like` operator", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddFilter("field-like=value")
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

	t.Run("Test AddFilter with invalid operator", func(t *testing.T) {
		qp := NewQueryParams()
		qp.AddFilter("field-invalid=value")
		if len(qp.Fitlers) != 0 {
			t.Errorf("Expected 0, got %d", len(qp.Fitlers))
		}
	})

	t.Run("Test Pagination", func(t *testing.T) {
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

	t.Run("Test Pagination with invalid offset or limit", func(t *testing.T) {
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

	t.Run("Test AddOrderBy", func(t *testing.T) {
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
}
