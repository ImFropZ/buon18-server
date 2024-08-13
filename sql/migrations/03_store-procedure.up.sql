CREATE
OR REPLACE PROCEDURE create_sales_order (
    name VARCHAR(64),
    commitment_date TIMESTAMP WITH TIME ZONE,
    note TEXT,
    sales_quotation_id BIGINT,
    accounting_payment_term_id BIGINT,
    cid BIGINT,
    ctime TIMESTAMP WITH TIME ZONE,
    mid BIGINT,
    mtime TIMESTAMP WITH TIME ZONE
) LANGUAGE plpgsql AS $$
DECLARE
    s_quotation_status sales_quotation_status_typ := 'quotation';

BEGIN

    -- Get sales quotation status
    SELECT
        "sales.quotation".status INTO s_quotation_status
    FROM
        "sales.quotation"
    WHERE
        "sales.quotation".id = sales_quotation_id;

    -- Create sales order if sales quotation is in sales order status
    IF s_quotation_status = 'sales_order' THEN
        INSERT INTO
            "sales.order"
            (name, commitment_date, note, sales_quotation_id, accounting_payment_term_id, cid, ctime, mid, mtime)
        VALUES
            (name, commitment_date, note, sales_quotation_id, accounting_payment_term_id, cid, ctime, mid, mtime);
    ELSE
        RAISE EXCEPTION 'custom_error:sales quotation is not in sales_order status';
    END IF;

END;

$$;