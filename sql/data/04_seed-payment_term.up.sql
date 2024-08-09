INSERT INTO
    "accounting.payment_term" (id, name, description, cid, ctime, mid, mtime)
VALUES
    (1, 'Net 30', 'Net 30', 1, NOW(), 1, NOW()),
    (2, 'Net 60', 'Net 60', 1, NOW(), 1, NOW()),
    (3, 'Net 90', 'Net 90', 1, NOW(), 1, NOW()),
    (
        4,
        '30% Now, Balance 60 Days',
        'Pay 30% now, balance due in 60 days',
        1,
        NOW(),
        1,
        NOW()
    );

INSERT INTO
    "accounting.payment_term_line" (
        id,
        sequence,
        value_amount_percent,
        number_of_days,
        accounting_payment_term_id,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (1, 1, 100, 30, 1, 1, NOW(), 1, NOW()),
    (2, 1, 100, 60, 2, 1, NOW(), 1, NOW()),
    (3, 1, 100, 90, 3, 1, NOW(), 1, NOW()),
    (4, 1, 30, 0, 4, 1, NOW(), 1, NOW()),
    (5, 2, 70, 60, 4, 1, NOW(), 1, NOW());