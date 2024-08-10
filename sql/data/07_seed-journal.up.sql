INSERT INTO
    "accounting.journal" (
        id,
        code,
        name,
        typ,
        accounting_account_id,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (
        1,
        'JNL1001',
        'Sales Journal',
        'sales',
        2,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        2,
        'JNL1002',
        'Purchase Journal',
        'purchase',
        4,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        3,
        'JNL1003',
        'Cash Journal',
        'cash',
        1,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        4,
        'JNL1004',
        'Bank Journal',
        'bank',
        3,
        1,
        NOW(),
        1,
        NOW()
    );