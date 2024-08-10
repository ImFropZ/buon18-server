INSERT INTO
    "sales.quotation" (
        id,
        name,
        creation_date,
        validity_date,
        discount,
        amount_delivery,
        status,
        setting_customer_id,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (
        1,
        'Quotation 1',
        '2021-01-01',
        '2021-01-31',
        50,
        100,
        'quotation',
        500,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        2,
        'Quotation 2',
        '2021-02-01',
        '2021-02-28',
        100,
        200,
        'quotation_sent',
        500,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        3,
        'Quotation 3',
        '2021-03-01',
        '2021-03-31',
        150,
        300,
        'quotation_sent',
        500,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        4,
        'Quotation 4',
        '2021-04-01',
        '2021-04-30',
        200,
        400,
        'sales_order',
        501,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        5,
        'Quotation 5',
        '2021-05-01',
        '2021-05-31',
        250,
        500,
        'sales_order',
        501,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        6,
        'Quotation 6',
        '2021-06-01',
        '2021-06-30',
        300,
        600,
        'cancelled',
        501,
        1,
        NOW(),
        1,
        NOW()
    );

INSERT INTO
    "sales.order_item" (
        id,
        name,
        description,
        price,
        discount,
        sales_quotation_id,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (
        1,
        'Item 1',
        'Item 1 description',
        100,
        0,
        1,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        2,
        'Item 2',
        'Item 2 description',
        200,
        0,
        1,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        3,
        'Item 3',
        'Item 3 description',
        500,
        50,
        2,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        4,
        'Item 4',
        'Item 4 description',
        1000,
        100,
        3,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        5,
        'Item 5',
        'Item 5 description',
        2000,
        150,
        4,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        6,
        'Item 6',
        'Item 6 description',
        3000,
        200,
        5,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        7,
        'Item 7',
        'Item 7 description',
        4000,
        250,
        6,
        1,
        NOW(),
        1,
        NOW()
    );