package database

import "github.com/lib/pq"

type PQ_KEY_CONSTRAINTS int

const (
	DUPLICATE PQ_KEY_CONSTRAINTS = iota
	VIOLATE_FOREIGN_KEY
	VIOLATE_CHECK
)

var PQ_ERROR_CODES = map[PQ_KEY_CONSTRAINTS]pq.Error{
	DUPLICATE: {
		Code: "23505",
	},
	VIOLATE_FOREIGN_KEY: {
		Code: "23503",
	},
	VIOLATE_CHECK: {
		Code: "23514",
	},
}

const (
	KEY_SETTING_USER_EMAIL = "setting.user_email_key"

	CHK_EXPIRY_DATE   = "chk_expiry_date"
	CHK_DELIVERY_DATE = "chk_delivery_date"
)
