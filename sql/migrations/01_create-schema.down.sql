-- DROP TABLE
DROP TABLE IF EXISTS "accounting.journal_entry_line";

DROP TABLE IF EXISTS "accounting.journal_entry";

DROP TABLE IF EXISTS "accounting.journal";

DROP TABLE IF EXISTS "accounting.account";

DROP TABLE IF EXISTS "sales.order";

DROP TABLE IF EXISTS "accounting.payment_term_line";

DROP TABLE IF EXISTS "accounting.payment_term";

DROP TABLE IF EXISTS "sales.order_item";

DROP TABLE IF EXISTS "sales.quotation";

DROP TABLE IF EXISTS "setting.customer";

DROP TABLE IF EXISTS "setting.user";

DROP TABLE IF EXISTS "setting.role_permission";

DROP TABLE IF EXISTS "setting.permission";

DROP TABLE IF EXISTS "setting.role";

-- DROP TYPE
DROP TYPE IF EXISTS "accounting_account_typ";

DROP TYPE IF EXISTS "accounting_journal_typ";

DROP TYPE IF EXISTS "accounting_journal_typaccounting_journal_entry_status_typ";

DROP TYPE IF EXISTS "sales_quotation_status_typ";

DROP TYPE IF EXISTS "setting_gender_typ";

DROP TYPE IF EXISTS "setting_user_typ";