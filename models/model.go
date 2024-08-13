package models

import "time"

var (
	// Gender
	SettingGenderTypMale   = "m"
	SettingGenderTypFemale = "f"
	SettingGenderTypOther  = "o"

	// Setting user types
	SettingUserTypUser = "user"
	SettingUserTypBot  = "bot"

	// Sales quotation status
	SalesQuotationStatusQuotation      = "quotation"
	SalesQuotationStatusQuotationSent  = "quotation_sent"
	SalesQuotationStatusSalesOrder     = "sales_order"
	SalesQuotationStatusSalesCancelled = "cancelled"

	// Accounting account types
	AccountingAccountTypAssetCurrent        = "asset_current"
	AccountingAccountTypAssetNonCurrent     = "asset_non_current"
	AccountingAccountTypLiabilityCurrent    = "liability_current"
	AccountingAccountTypLiabilityNonCurrent = "liability_non_current"
	AccountingAccountTypEquity              = "equity"
	AccountingAccountTypIncome              = "income"
	AccountingAccountTypExpense             = "expense"
	AccountingAccountTypGain                = "gain"
	AccountingAccountTypLoss                = "loss"
)

var VALID_GENDER_TYPES = []string{SettingGenderTypMale, SettingGenderTypFemale, SettingGenderTypOther}
var VALID_SALES_QUOTATION_STATUS = []string{SalesQuotationStatusQuotation, SalesQuotationStatusQuotationSent, SalesQuotationStatusSalesOrder, SalesQuotationStatusSalesCancelled}
var VALID_ACCOUNTING_ACCOUNT_TYPES = []string{AccountingAccountTypAssetCurrent, AccountingAccountTypAssetNonCurrent, AccountingAccountTypLiabilityCurrent, AccountingAccountTypLiabilityNonCurrent, AccountingAccountTypEquity, AccountingAccountTypIncome, AccountingAccountTypExpense, AccountingAccountTypGain, AccountingAccountTypLoss}

type CommonModel struct {
	CId   uint
	CTime time.Time
	MId   uint
	MTime time.Time
}

func (cm *CommonModel) PrepareForCreate(cid uint, mid uint) (err error) {
	cm.CId = cid
	cm.CTime = time.Now()
	cm.MId = mid
	cm.MTime = time.Now()
	return
}

func (cm *CommonModel) PrepareForUpdate(mid uint) (err error) {
	cm.MId = mid
	cm.MTime = time.Now()
	return
}
