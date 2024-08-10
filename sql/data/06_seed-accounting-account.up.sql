INSERT INTO
    "accounting.account" (id, code, name, typ, cid, ctime, mid, mtime)
VALUES
    (
        1,
        'AC1001',
        'Cash',
        'asset_current',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        2,
        'AC1002',
        'Accounts Receivable',
        'asset_non_current',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        3,
        'LC1003',
        'Accounts Payable',
        'liability_current',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        4,
        'LC1004',
        'Long-term Debt',
        'liability_non_current',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        5,
        'EQ1005',
        'Common Stock',
        'equity',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        6,
        'IN1006',
        'Sales Revenue',
        'income',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        7,
        'EX1007',
        'Rent Expense',
        'expense',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        8,
        'GN1008',
        'Gain on Sale',
        'gain',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        9,
        'LS1009',
        'Loss on Sale',
        'loss',
        1,
        NOW(),
        1,
        NOW()
    );