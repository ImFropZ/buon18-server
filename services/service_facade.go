package services

import (
	"server/services/accounting"
	"server/services/sales"
	"server/services/setting"
)

type ServiceFacade struct {
	AuthService                   *AuthService
	SettingCustomerService        *setting.SettingCustomerService
	SettingRoleService            *setting.SettingRoleService
	SettingUserService            *setting.SettingUserService
	SettingPermissionService      *setting.SettingPermissionService
	SalesOrderService             *sales.SalesOrderService
	SalesQuotationService         *sales.SalesQuotationService
	AccountingAccountService      *accounting.AccountingAccountService
	AccountingJournalEntryService *accounting.AccountingJournalEntryService
	AccountingJournalService      *accounting.AccountingJournalService
	AccountingPaymentTermService  *accounting.AccountingPaymentTermService
}
