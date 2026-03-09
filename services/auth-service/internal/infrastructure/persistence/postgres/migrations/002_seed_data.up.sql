SET search_path TO imcs_auth;

INSERT INTO permissions (code, name, module) VALUES
    ('users.list',   'List Users',    'users'),
    ('users.create', 'Create User',   'users'),
    ('users.update', 'Update User',   'users'),
    ('users.delete', 'Delete User',   'users'),
    ('roles.manage', 'Manage Roles',  'roles'),
    ('config.view',  'View Config',   'config'),
    ('config.manage','Manage Config',  'config'),
    ('franchises.view',   'View Franchises',   'franchises'),
    ('franchises.create', 'Create Franchise',  'franchises'),
    ('franchises.update', 'Update Franchise',  'franchises'),
    ('shipments.view',   'View Shipments',   'shipments'),
    ('shipments.create', 'Create Shipment',  'shipments'),
    ('shipments.manage', 'Manage Shipments', 'shipments'),
    ('customers.view',   'View Customers',   'customers'),
    ('customers.create', 'Create Customer',  'customers'),
    ('customers.manage', 'Manage Customers', 'customers'),
    ('invoices.view',    'View Invoices',    'invoices'),
    ('invoices.manage',  'Manage Invoices',  'invoices'),
    ('reports.view',     'View Reports',     'reports')
ON CONFLICT (code) DO NOTHING;

INSERT INTO roles (country_code, name, display_name, description, is_system) VALUES
    ('AU', 'super_admin',    'Super Administrator',     'Full system access across all countries', true),
    ('AU', 'admin',          'Administrator',           'Full access for country',                true),
    ('AU', 'franchise_admin','Franchise Administrator', 'Franchise management access',            true),
    ('AU', 'franchise_user', 'Franchise User',          'Day-to-day franchise operations',        true),
    ('AU', 'customer',       'Customer',                'Customer portal access',                 true)
ON CONFLICT (country_code, name) DO NOTHING;

-- Assign all permissions to admin
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r CROSS JOIN permissions p
WHERE r.name = 'admin' AND r.country_code = 'AU'
ON CONFLICT DO NOTHING;
