-- Check date constraint
ALTER TABLE "quote"
DROP CONSTRAINT chk_expiry_date;

ALTER TABLE "sales_order"
DROP CONSTRAINT chk_delivery_date;