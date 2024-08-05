package database

type PQ_KEY_CONSTRAINTS int

const (
	DUPLICATE PQ_KEY_CONSTRAINTS = iota
	VIOLATE_FOREIGN_KEY
	VIOLATE_CHECK
)

var PQ_ERROR_CODES = map[PQ_KEY_CONSTRAINTS]string{
	DUPLICATE:           "23505",
	VIOLATE_FOREIGN_KEY: "23503",
	VIOLATE_CHECK:       "23514",
}

const (
	CHK_EXPIRY_DATE   = "chk_expiry_date"
	CHK_DELIVERY_DATE = "chk_delivery_date"
)
