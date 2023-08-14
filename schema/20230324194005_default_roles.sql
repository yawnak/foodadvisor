-- Create default roles
INSERT INTO roles (name)
VALUES ('user'),
    ('admin'),
    ('owner');
-- Add default permissions for admin
INSERT INTO permissions_to_roles (role, permission)
VALUES ('admin', 'editUserRole');
-- Add all permissions for owner
INSERT INTO permissions_to_roles (role, permission)
SELECT roles.name AS role,
    permissions.name AS permission
FROM roles,
    permissions
WHERE roles.name = 'owner';