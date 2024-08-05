-- Check date constraint
ALTER TABLE "quote"
ADD CONSTRAINT chk_expiry_date CHECK (expiry_date >= date);

ALTER TABLE "sales_order"
ADD CONSTRAINT chk_delivery_date CHECK (delivery_date >= accept_date);