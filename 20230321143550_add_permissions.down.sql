DELETE FROM permissions
WHERE name IN (
    "editRoles",
    "editUserRole"
);