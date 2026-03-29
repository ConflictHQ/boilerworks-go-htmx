-- Combined init script for Docker (runs both schema and seed)

-- Users
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sessions
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- Groups
CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Permissions
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Group permissions
CREATE TABLE IF NOT EXISTS group_permissions (
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, permission_id)
);

-- User groups
CREATE TABLE IF NOT EXISTS user_groups (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

-- Categories
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories(deleted_at);

-- Items
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price DECIMAL(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    category_id UUID REFERENCES categories(id),
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_items_deleted_at ON items(deleted_at);
CREATE INDEX IF NOT EXISTS idx_items_status ON items(status);
CREATE INDEX IF NOT EXISTS idx_items_category_id ON items(category_id);

-- Form definitions
CREATE TABLE IF NOT EXISTS form_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    schema JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_form_definitions_slug ON form_definitions(slug);
CREATE INDEX IF NOT EXISTS idx_form_definitions_deleted_at ON form_definitions(deleted_at);

-- Form submissions
CREATE TABLE IF NOT EXISTS form_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    form_definition_id UUID NOT NULL REFERENCES form_definitions(id),
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_form_submissions_form_definition_id ON form_submissions(form_definition_id);

-- Workflow definitions
CREATE TABLE IF NOT EXISTS workflow_definitions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status VARCHAR(50) NOT NULL DEFAULT 'draft',
    states JSONB NOT NULL DEFAULT '[]'::jsonb,
    transitions JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_by UUID NOT NULL REFERENCES users(id),
    updated_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_workflow_definitions_deleted_at ON workflow_definitions(deleted_at);

-- Workflow instances
CREATE TABLE IF NOT EXISTS workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    workflow_definition_id UUID NOT NULL REFERENCES workflow_definitions(id),
    current_state VARCHAR(255) NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workflow_instances_definition_id ON workflow_instances(workflow_definition_id);

-- Transition logs
CREATE TABLE IF NOT EXISTS transition_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_instance_id UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    from_state VARCHAR(255) NOT NULL,
    to_state VARCHAR(255) NOT NULL,
    performed_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transition_logs_instance_id ON transition_logs(workflow_instance_id);

-- Seed data
INSERT INTO permissions (name) VALUES
    ('items.view'), ('items.create'), ('items.edit'), ('items.delete'),
    ('categories.view'), ('categories.create'), ('categories.edit'), ('categories.delete'),
    ('forms.view'), ('forms.create'), ('forms.edit'), ('forms.delete'),
    ('workflows.view'), ('workflows.create'), ('workflows.edit'), ('workflows.delete'),
    ('users.view'), ('users.create'), ('users.edit'), ('users.delete')
ON CONFLICT (name) DO NOTHING;

INSERT INTO groups (name) VALUES ('admin'), ('editor'), ('viewer')
ON CONFLICT (name) DO NOTHING;

-- Admin gets all permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p WHERE g.name = 'admin'
ON CONFLICT DO NOTHING;

-- Editor gets create/edit/view permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p
WHERE g.name = 'editor' AND (p.name LIKE '%.create' OR p.name LIKE '%.edit' OR p.name LIKE '%.view')
ON CONFLICT DO NOTHING;

-- Viewer gets view-only permissions
INSERT INTO group_permissions (group_id, permission_id)
SELECT g.id, p.id FROM groups g, permissions p
WHERE g.name = 'viewer' AND p.name LIKE '%.view'
ON CONFLICT DO NOTHING;

-- Admin user (password: "password")
INSERT INTO users (name, email, password_hash) VALUES
    ('Admin', 'admin@boilerworks.dev', '$2a$10$6Fj7ifteQL4P3jEmgmeYbeukJ0QMNGVDjsL0c5MX/SsTixnH9c7Oe')
ON CONFLICT (email) DO NOTHING;

-- Add admin to admin group
INSERT INTO user_groups (user_id, group_id)
SELECT u.id, g.id FROM users u, groups g WHERE u.email = 'admin@boilerworks.dev' AND g.name = 'admin'
ON CONFLICT DO NOTHING;
