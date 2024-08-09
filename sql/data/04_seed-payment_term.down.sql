DELETE FROM "accounting.payment_term_line"
WHERE
    id < 1000;

DELETE FROM "accounting.payment_term"
WHERE
    id < 1000;