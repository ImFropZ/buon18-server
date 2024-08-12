INSERT INTO
    "setting.role" (id, name, description, cid, ctime, mid, mtime)
VALUES
    (
        3,
        'Setting Administrator',
        'Full access to all settings',
        1,
        NOW(),
        1,
        NOW()
    ), -- Admin permissions
    (
        4,
        'Sales Administrator',
        'Full access to all sales',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        5,
        'Accounting Administrator',
        'Full access to all accounting',
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
    (3, 3, 1, NOW(), 1, NOW()), -- Admin permissions
    (4, 4, 1, NOW(), 1, NOW()),
    (5, 5, 1, NOW(), 1, NOW());

-- User permissions
INSERT INTO
    "setting.user" (
        id,
        name,
        email,
        setting_role_id,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (
        3,
        'Setting Admin',
        'setting@buon18.com',
        3,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        4,
        'Sales Admin',
        'sales@buon18.com',
        4,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        5,
        'Accounting Admin',
        'accounting@buon18.com',
        5,
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        1000,
        'John Doe',
        'jd@dummy-data.com',
        2,
        1,
        NOW(),
        1,
        NOW()
    ), -- User permissions
    (
        1001,
        'Jane Doe',
        'jad@dummy-data.com',
        2,
        1,
        NOW(),
        1,
        NOW()
    );