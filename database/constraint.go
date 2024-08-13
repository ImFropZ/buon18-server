package database

const (
	KEY_SETTING_USER_EMAIL      = "setting.user_email_key"
	KEY_SETTING_CUSTOMER_EMAIL  = "setting.customer_email_key"
	KEY_SALES_QUOTATION_NAME    = "sales.quotation_name_key"
	KEY_SALES_ORDER_NAME        = "sales.order_name_key"
	KEY_ACCOUNTING_ACCOUNT_CODE = "accounting.account_code_key"

	FK_SALES_QUOTATION_CUSTOMER_ID = "setting.customer_id_fkey"
	FK_ACCOUNTING_PAYMENT_TERM_ID  = "accounting.payment_term_id_fkey"

	CHK_SALES_QUOTATION_DATE = "sales.quotation_date_chk"
)
