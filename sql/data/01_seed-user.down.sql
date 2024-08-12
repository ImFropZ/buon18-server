DELETE FROM "setting.user"
WHERE
    email ILIKE '%@dummy-data.com'
    OR id BETWEEN 3 AND 5;

DELETE FROM "setting.role_permission"
WHERE
    setting_role_id BETWEEN 3 AND 5;

DELETE FROM "setting.role"
WHERE
    id BETWEEN 3 AND 5;