CREATE
OR REPLACE PROCEDURE update_quote_total(quote_id INT) LANGUAGE plpgsql AS $$
DECLARE
    v_subtotal NUMERIC := 0;

BEGIN
    -- Calculate the subtotal for the specified qid
    SELECT
        COALESCE(SUM(qi.quantity * qi.unit_price), 0) INTO v_subtotal
    FROM
        "quote_item" AS qi
    WHERE
        qi.quote_id = qid;

-- Update the quote with the calculated subtotal and total
UPDATE
    "quote"
SET
    subtotal = v_subtotal,
    total = v_subtotal - discount
WHERE
    id = qid;

END;

$$;