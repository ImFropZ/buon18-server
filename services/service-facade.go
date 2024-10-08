package services

import (
	"system.buon18.com/m/services/accounting"
	"system.buon18.com/m/services/sales"
	"system.buon18.com/m/services/setting"
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
