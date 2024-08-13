#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    INSERT INTO
        "sales.order" (
            id,
            name,
            commitment_date,
            note,
            sales_quotation_id,
            accounting_payment_term_id,
            cid,
            ctime,
            mid,
            mtime
        )
    VALUES
        (
            1,
            'Order 1',
            '2021-04-05',
            '',
            4,
            4,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            2,
            'Order 2',
            '2021-05-05',
            '',
            5,
            2,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            3,
            'Order 3',
            '2021-06-05',
            '',
            6,
            1,
            1,
            NOW(),
            1,
            NOW()
        );
EOSQL
