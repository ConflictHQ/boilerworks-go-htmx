-- +goose Up

-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Sessions
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- Groups
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Permissions
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Group permissions (many-to-many)
CREATE TABLE group_permissions (
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, permission_id)
);

-- User groups (many-to-many)
CREATE TABLE user_groups (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);

-- Categories
CREATE TABLE categories (
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

CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);

-- Products
CREATE TABLE products (
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

CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE INDEX idx_products_status ON products(status);
CREATE INDEX idx_products_category_id ON products(category_id);

-- Form definitions
CREATE TABLE form_definitions (
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

CREATE INDEX idx_form_definitions_slug ON form_definitions(slug);
CREATE INDEX idx_form_definitions_deleted_at ON form_definitions(deleted_at);

-- Form submissions
CREATE TABLE form_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    form_definition_id UUID NOT NULL REFERENCES form_definitions(id),
    data JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_form_submissions_form_definition_id ON form_submissions(form_definition_id);

-- Workflow definitions
CREATE TABLE workflow_definitions (
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

CREATE INDEX idx_workflow_definitions_deleted_at ON workflow_definitions(deleted_at);

-- Workflow instances
CREATE TABLE workflow_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uuid UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    workflow_definition_id UUID NOT NULL REFERENCES workflow_definitions(id),
    current_state VARCHAR(255) NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workflow_instances_definition_id ON workflow_instances(workflow_definition_id);

-- Transition logs
CREATE TABLE transition_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_instance_id UUID NOT NULL REFERENCES workflow_instances(id) ON DELETE CASCADE,
    from_state VARCHAR(255) NOT NULL,
    to_state VARCHAR(255) NOT NULL,
    performed_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transition_logs_instance_id ON transition_logs(workflow_instance_id);

-- +goose Down
DROP TABLE IF EXISTS transition_logs;
DROP TABLE IF EXISTS workflow_instances;
DROP TABLE IF EXISTS workflow_definitions;
DROP TABLE IF EXISTS form_submissions;
DROP TABLE IF EXISTS form_definitions;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS user_groups;
DROP TABLE IF EXISTS group_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
