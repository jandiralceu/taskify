-- Create roles table
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create permissions table
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(resource, action)
);

-- Create user_roles table
CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Create role_permissions table
CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Casbin rules table (required for GORM adapter)
CREATE TABLE casbin_rule (
    id BIGSERIAL PRIMARY KEY,
    ptype VARCHAR(100),
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100),
    v3 VARCHAR(100),
    v4 VARCHAR(100),
    v5 VARCHAR(100)
);

CREATE INDEX idx_casbin_rule ON casbin_rule (ptype, v0, v1, v2, v3, v4, v5);

-- Seed initial roles
INSERT INTO roles (name, description) VALUES 
('admin', 'Full system access'),
('employee', 'Standard user access');

-- Seed standard permissions
INSERT INTO permissions (resource, action, description) VALUES 
('*', '*', 'All actions on all resources'),
('/api/v1/tasks*', '*', 'All actions on tasks'),
('/api/v1/users/profile', 'GET', 'View own profile'),
('/api/v1/users/profile', 'PATCH', 'Update own profile'),
('/api/v1/users/change-password', 'PATCH', 'Change own password'),
('/api/v1/users/avatar', 'POST', 'Update own avatar');

-- Map permissions to roles
-- Admin gets everything
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'admin' AND p.resource = '*' AND p.action = '*';

-- Employee gets restricted access
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'employee' AND (
    (p.resource = '/api/v1/tasks*' AND p.action = '*') OR
    (p.resource = '/api/v1/users/profile' AND p.action = 'GET') OR
    (p.resource = '/api/v1/users/profile' AND p.action = 'PATCH') OR
    (p.resource = '/api/v1/users/change-password' AND p.action = 'PATCH') OR
    (p.resource = '/api/v1/users/avatar' AND p.action = 'POST')
);

-- Assign existing users to 'employee' role by default (migration)
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id FROM users u, roles r 
WHERE r.name = 'employee';

-- Seed initial Casbin rules (for GORM adapter)
INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES
('p', 'admin', '/api/v1/*', '*'),
('p', 'employee', '/api/v1/auth/signout', 'POST'),
('p', 'employee', '/api/v1/auth/refresh', 'POST'),
('p', 'employee', '/api/v1/users/permissions', 'GET'),
('p', 'employee', '/api/v1/users/profile', 'GET'),
('p', 'employee', '/api/v1/users/profile', 'PATCH'),
('p', 'employee', '/api/v1/users/profile', 'DELETE'),
('p', 'employee', '/api/v1/users/change-password', 'PATCH'),
('p', 'employee', '/api/v1/users/avatar', 'POST'),
('p', 'employee', '/api/v1/tasks*', '*'),
('p', 'employee', '/api/v1/tasks/*', '*');
