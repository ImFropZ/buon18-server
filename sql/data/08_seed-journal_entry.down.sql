DELETE FROM "accounting.journal_entry_line"
WHERE
    id BETWEEN 1 AND 7;

DELETE FROM "accounting.journal_entry"
WHERE
    id BETWEEN 1 AND 4;