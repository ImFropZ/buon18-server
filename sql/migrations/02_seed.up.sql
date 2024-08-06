INSERT INTO
    "setting.role" (id, name, description, cid, ctime, mid, mtime)
VALUES
    (1, 'bot', 'BOT', 1, NOW(), 1, NOW());

INSERT INTO
    "setting.permission" (id, name, cid, ctime, mid, mtime)
VALUES
    (1, 'full access', 1, NOW(), 1, NOW());

INSERT INTO
    "setting.role_permission" (setting_role_id, setting_permission_id, cid, ctime, mid, mtime)
VALUES
    (1, 1, 1, NOW(), 1, NOW());

INSERT INTO
    "setting.user" (
        id,
        name,
        pwd,
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
        '',
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
        '',
        'admin@buon18.com',
        'user',
        1,
        1,
        NOW(),
        1,
        NOW()
    );