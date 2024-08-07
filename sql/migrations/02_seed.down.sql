DELETE FROM "setting.user"
WHERE
    id < 1000;

DELETE FROM "setting.role_permission"
WHERE
    setting_role_id < 1000;

DELETE FROM "setting.permission"
WHERE
    id < 1000;

DELETE FROM "setting.role"
WHERE
    id < 1000;