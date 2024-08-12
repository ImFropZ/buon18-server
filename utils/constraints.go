package utils

var (
	FULL_ACCESS_ID     = 1
	FULL_AUTH_ID       = 2
	FULL_SETTING_ID    = 3
	FULL_SALES_ID      = 4
	FULL_ACCOUNTING_ID = 5
)

var FULL_PERMISSION_IDS = []string{
	IntToStr(FULL_ACCESS_ID),
	IntToStr(FULL_AUTH_ID),
	IntToStr(FULL_SETTING_ID),
	IntToStr(FULL_SALES_ID),
	IntToStr(FULL_ACCOUNTING_ID),
}

type CommonPermission struct {
	VIEW   string
	CREATE string
	UPDATE string
	DELETE string
}

type Permissions struct {
	FULL_ACCESS                string
	FULL_AUTH                  string
	FULL_SETTING               string
	FULL_SALES                 string
	FULL_ACCOUNTING            string
	AUTH                       CommonPermission
	SETTING_USERS              CommonPermission
	SETTING_CUSTOMERS          CommonPermission
	SETTING_ROLES              CommonPermission
	SALES_QUOTATIONS           CommonPermission
	SALES_ORDERS               CommonPermission
	ACCOUNTING_ACCOUNTS        CommonPermission
	ACCOUNTING_JOURNALS        CommonPermission
	ACCOUNTING_JOURNAL_ENTRIES CommonPermission
	ACCOUNTING_PAYMENT_TERMS   CommonPermission
}

var PREDEFINED_PERMISSIONS = Permissions{
	FULL_ACCESS:     "FULL_ACCESS",
	FULL_AUTH:       "FULL_AUTH",
	FULL_SETTING:    "FULL_SETTING",
	FULL_SALES:      "FULL_SALES",
	FULL_ACCOUNTING: "FULL_ACCOUNTING",
	AUTH: CommonPermission{
		VIEW:   "VIEW_PROFILE",
		UPDATE: "UPDATE_PROFILE",
	},
	SETTING_USERS: CommonPermission{
		VIEW:   "VIEW_SETTING_USERS",
		CREATE: "CREATE_SETTING_USERS",
		UPDATE: "UPDATE_SETTING_USERS",
		DELETE: "DELETE_SETTING_USERS",
	},
	SETTING_CUSTOMERS: CommonPermission{
		VIEW:   "VIEW_SETTING_CUSTOMERS",
		CREATE: "CREATE_SETTING_CUSTOMERS",
		UPDATE: "UPDATE_SETTING_CUSTOMERS",
		DELETE: "DELETE_SETTING_CUSTOMERS",
	},
	SETTING_ROLES: CommonPermission{
		VIEW:   "VIEW_SETTING_ROLES",
		CREATE: "CREATE_SETTING_ROLES",
		UPDATE: "UPDATE_SETTING_ROLES",
		DELETE: "DELETE_SETTING_ROLES",
	},
	SALES_QUOTATIONS: CommonPermission{
		VIEW:   "VIEW_SALES_QUOTATIONS",
		CREATE: "CREATE_SALES_QUOTATIONS",
		UPDATE: "UPDATE_SALES_QUOTATIONS",
		DELETE: "DELETE_SALES_QUOTATIONS",
	},
	SALES_ORDERS: CommonPermission{
		VIEW:   "VIEW_SALES_ORDERS",
		CREATE: "CREATE_SALES_ORDERS",
		UPDATE: "UPDATE_SALES_ORDERS",
		DELETE: "DELETE_SALES_ORDERS",
	},
	ACCOUNTING_ACCOUNTS: CommonPermission{
		VIEW:   "VIEW_ACCOUNTING_ACCOUNTS",
		CREATE: "CREATE_ACCOUNTING_ACCOUNTS",
		UPDATE: "UPDATE_ACCOUNTING_ACCOUNTS",
		DELETE: "DELETE_ACCOUNTING_ACCOUNTS",
	},
	ACCOUNTING_JOURNALS: CommonPermission{
		VIEW:   "VIEW_ACCOUNTING_JOURNALS",
		CREATE: "CREATE_ACCOUNTING_JOURNALS",
		UPDATE: "UPDATE_ACCOUNTING_JOURNALS",
		DELETE: "DELETE_ACCOUNTING_JOURNALS",
	},
	ACCOUNTING_JOURNAL_ENTRIES: CommonPermission{
		VIEW:   "VIEW_ACCOUNTING_JOURNAL_ENTRIES",
		CREATE: "CREATE_ACCOUNTING_JOURNAL_ENTRIES",
		UPDATE: "UPDATE_ACCOUNTING_JOURNAL_ENTRIES",
		DELETE: "DELETE_ACCOUNTING_JOURNAL_ENTRIES",
	},
	ACCOUNTING_PAYMENT_TERMS: CommonPermission{
		VIEW:   "VIEW_ACCOUNTING_PAYMENT_TERMS",
		CREATE: "CREATE_ACCOUNTING_PAYMENT_TERMS",
		UPDATE: "UPDATE_ACCOUNTING_PAYMENT_TERMS",
		DELETE: "DELETE_ACCOUNTING_PAYMENT_TERMS",
	},
}
