-- +goose Up

-- Seed permissions
INSERT INTO permissions (name) VALUES
    ('items.view'), ('items.create'), ('items.edit'), ('items.delete'),
    ('categories.view'), ('categories.create'), ('categories.edit'), ('categories.delete'),
    ('forms.view'), ('forms.create'), ('forms.edit'), ('forms.delete'),
    ('workflows.view'), ('workflows.create'), ('workflows.edit'), ('workflows.delete'),
    ('users.view'), ('users.create'), ('users.edit'), ('users.delete');

-- Seed groups
INSERT INTO groups (name) VALUES ('admin'), ('editor'), ('viewer');

-- Admin gets all permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p WHERE g.name = 'admin';

-- Editor gets create/edit/view permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p
WHERE g.name = 'editor' AND (p.name LIKE '%.create' OR p.name LIKE '%.edit' OR p.name LIKE '%.view');

-- Viewer gets view-only permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p
WHERE g.name = 'viewer' AND p.name LIKE '%.view';

-- Seed admin user (password: "password")
-- bcrypt hash of "password" with cost 10
INSERT INTO users (name, email, password_hash) VALUES
    ('Admin', 'admin@boilerworks.dev', '$2a$10$6Fj7ifteQL4P3jEmgmeYbeukJ0QMNGVDjsL0c5MX/SsTixnH9c7Oe');

-- Add admin to admin group
INSERT INTO user_groups (user_id, group_id)
SELECT u.id, g.id FROM users u, groups g WHERE u.email = 'admin@boilerworks.dev' AND g.name = 'admin';

-- +goose Down
DELETE FROM user_groups;
DELETE FROM group_permissions;
DELETE FROM users WHERE email = 'admin@boilerworks.dev';
DELETE FROM groups;
DELETE FROM permissions;
