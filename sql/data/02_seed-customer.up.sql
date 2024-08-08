INSERT INTO
    "setting.customer" (
        id,
        fullname,
        gender,
        email,
        phone,
        additional_information,
        cid,
        ctime,
        mid,
        mtime
    )
VALUES
    (
        500,
        'John Doe',
        'm',
        'jd@dummy-data.com',
        '096123456',
        '{"note":"This is a dummy data from john doe"}',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        501,
        'Jane Doe',
        'f',
        'jad@dummy-data.com',
        '064456789',
        '{"note":"This is a dummy data from jane doe"}',
        1,
        NOW(),
        1,
        NOW()
    ),
    (
        502,
        'John Foo',
        'u',
        'jf@dummy-data.com',
        '012789123',
        '{"note":"This is a dummy data from john foo"}',
        1,
        NOW(),
        1,
        NOW()
    );