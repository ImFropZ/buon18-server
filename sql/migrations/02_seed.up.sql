INSERT INTO
    "setting.role" (id, name, description, cid, ctime, mid, mtime)
VALUES
    (1, 'bot', 'BOT', 1, NOW(), 1, NOW()),
    (2, 'user', 'User', 1, NOW(), 1, NOW());

INSERT INTO
    "setting.permission" (id, name, cid, ctime, mid, mtime)
VALUES
    (1, 'FULL_ACESS', 1, NOW(), 1, NOW()), -- Permissions
    (2, 'FULL_AUTH', 1, NOW(), 1, NOW()),
    (3, 'FULL_SETTING', 1, NOW(), 1, NOW()),
    (4, 'FULL_SALES', 1, NOW(), 1, NOW()),
    (5, 'FULL_ACCOUNTING', 1, NOW(), 1, NOW()),
    (6, 'VIEW_PROFILE', 1, NOW(), 1, NOW()), -- Auth Permissions
    (7, 'UPDATE_PROFILE', 1, NOW(), 1, NOW()),
    (8, 'VIEW_SETTING_USERS', 1, NOW(), 1, NOW()), -- Setting Permissions
    (9, 'CREATE_SETTING_USERS', 1, NOW(), 1, NOW()),
    (10, 'UPDATE_SETTING_USERS', 1, NOW(), 1, NOW()),
    (11, 'DELETE_SETTING_USERS', 1, NOW(), 1, NOW()),
    (12, 'VIEW_SETTING_CUSTOMERS', 1, NOW(), 1, NOW()),
    (
        13,
        'CREATE_SETTING_CUSTOMERS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        14,
        'UPDATE_SETTING_CUSTOMERS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        15,
        'DELETE_SETTING_CUSTOMERS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (16, 'VIEW_SETTING_ROLES', 1, NOW(), 1, NOW()),
    (17, 'CREATE_SETTING_ROLES', 1, NOW(), 1, NOW()),
    (18, 'UPDATE_SETTING_ROLES', 1, NOW(), 1, NOW()),
    (19, 'DELETE_SETTING_ROLES', 1, NOW(), 1, NOW()),
    (20, 'VIEW_SALES_QUOTATIONS', 1, NOW(), 1, NOW()), -- Sales Permissions
    (21, 'CREATE_SALES_QUOTATIONS', 1, NOW(), 1, NOW()),
    (22, 'UPDATE_SALES_QUOTATIONS', 1, NOW(), 1, NOW()),
    (23, 'DELETE_SALES_QUOTATIONS', 1, NOW(), 1, NOW()),
    (24, 'VIEW_SALES_ORDERS', 1, NOW(), 1, NOW()),
    (25, 'CREATE_SALES_ORDERS', 1, NOW(), 1, NOW()),
    (26, 'UPDATE_SALES_ORDERS', 1, NOW(), 1, NOW()),
    (27, 'DELETE_SALES_ORDERS', 1, NOW(), 1, NOW()),
    (
        28,
        'VIEW_ACCOUNTING_ACCOUNTS',
        1,
        NOW(),
        1,
        NOW()
    ), -- Accounting Permissions
    (
        29,
        'CREATE_ACCOUNTING_ACCOUNTS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        30,
        'UPDATE_ACCOUNTING_ACCOUNTS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        31,
        'DELETE_ACCOUNTING_ACCOUNTS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        32,
        'VIEW_ACCOUNTING_JOURNAL_ENTRIES',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        33,
        'CREATE_ACCOUNTING_JOURNAL_ENTRIES',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        34,
        'UPDATE_ACCOUNTING_JOURNAL_ENTRIES',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        35,
        'DELETE_ACCOUNTING_JOURNAL_ENTRIES',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        36,
        'VIEW_ACCOUNTING_PAYMENT_TERMS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        37,
        'CREATE_ACCOUNTING_PAYMENT_TERMS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        38,
        'UPDATE_ACCOUNTING_PAYMENT_TERMS',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        39,
        'DELETE_ACCOUNTING_PAYMENT_TERMS',
        1,
        NOW(),
        1,
        NOW()
    );

INSERT INTO
    "setting.role_permission" (
        setting_role_id,
        setting_permission_id,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (1, 1, 1, NOW(), 1, NOW()), -- Bot permissions
    (2, 6, 1, NOW(), 1, NOW()), -- User permissions
    (2, 7, 1, NOW(), 1, NOW());

INSERT INTO
    "setting.user" (
        id,
        name,
        email,
        typ,
        setting_role_id,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (
        1,
        'bot',
        'bot@buon18.com',
        'bot',
        1,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        2,
        'admin',
        'admin@buon18.com',
        'user',
        1,
        1,
        NOW(),
        1,
        NOW()
    );