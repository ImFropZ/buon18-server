#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    INSERT INTO
        "accounting.journal_entry" (
            id,
            name,
            date,
            note,
            status,
            accounting_journal_id,
            cid,
            ctime,
            mid,
            mtime
        )
    VALUES
        (
            1,
            'JE1001',
            '2024-08-09T00:00:00Z',
            'Entry for Sales Journal',
            'posted',
            1,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            2,
            'JE1002',
            '2024-08-10T00:00:00Z',
            'Entry for Purchase Journal',
            'draft',
            2,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            3,
            'JE1003',
            '2024-08-11T00:00:00Z',
            'Entry for Cash Journal',
            'posted',
            3,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            4,
            'JE1004',
            '2024-08-12T00:00:00Z',
            'Entry for Bank Journal',
            'cancelled',
            4,
            1,
            NOW(),
            1,
            NOW()
        );

    INSERT INTO
        "accounting.journal_entry_line" (
            id,
            sequence,
            name,
            amount_debit,
            amount_credit,
            accounting_journal_entry_id,
            accounting_account_id,
            cid,
            ctime,
            mid,
            mtime
        )
    VALUES
        (
            1,
            1,
            'Line 1 for JE1001',
            100.00,
            0.00,
            1,
            1,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            2,
            2,
            'Line 2 for JE1001',
            0.00,
            100.00,
            1,
            2,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            3,
            1,
            'Line 1 for JE1002',
            50.00,
            0.00,
            2,
            3,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            4,
            2,
            'Line 2 for JE1002',
            0.00,
            50.00,
            2,
            4,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            5,
            1,
            'Line 1 for JE1003',
            200.00,
            0.00,
            3,
            1,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            6,
            1,
            'Line 2 for JE1003',
            0.00,
            200.00,
            3,
            2,
            1,
            NOW(),
            1,
            NOW()
        ),
        (
            7,
            1,
            'Line 1 for JE1004',
            0.00,
            150.00,
            4,
            3,
            1,
            NOW(),
            1,
            NOW()
        );
EOSQL