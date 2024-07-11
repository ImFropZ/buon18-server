package database

type PQ_KEY_CONSTRAINTS int

const (
	DUPLICATE PQ_KEY_CONSTRAINTS = iota
	VIOLATE_FOREIGN_KEY
)

var PQ_ERROR_CODES = map[PQ_KEY_CONSTRAINTS]string{
	DUPLICATE:           "23505",
	VIOLATE_FOREIGN_KEY: "23503",
}
