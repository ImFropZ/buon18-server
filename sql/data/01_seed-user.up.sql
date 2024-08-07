INSERT INTO
    "setting.user" (
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
        'John Doe',
        'jd@dummy-data.com',
        2,
        1,
        NOW(),
        1,
        NOW()
    ), -- User permissions
    (
        'Jane Doe',
        'jad@dummy-data.com',
        2,
        1,
        NOW(),
        1,
        NOW()
    );