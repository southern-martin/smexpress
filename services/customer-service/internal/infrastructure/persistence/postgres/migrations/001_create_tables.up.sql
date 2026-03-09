SET search_path TO imcs_customers;

CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    franchise_id UUID,
    company_name VARCHAR(200) NOT NULL,
    trading_name VARCHAR(200),
    account_number VARCHAR(50),
    abn VARCHAR(20),
    email VARCHAR(200),
    phone VARCHAR(20),
    website VARCHAR(200),
    credit_limit NUMERIC(14,2) NOT NULL DEFAULT 0,
    credit_balance NUMERIC(14,2) NOT NULL DEFAULT 0,
    payment_terms INT NOT NULL DEFAULT 30,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_credit_hold BOOLEAN NOT NULL DEFAULT false,
    notes TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, account_number)
);

CREATE TABLE customer_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(200),
    phone VARCHAR(20),
    mobile VARCHAR(20),
    position VARCHAR(100),
    is_primary BOOLEAN NOT NULL DEFAULT false,
    is_billing BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE customer_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    address_type VARCHAR(20) NOT NULL DEFAULT 'shipping',
    company_name VARCHAR(200),
    contact_name VARCHAR(200),
    address_line1 VARCHAR(200) NOT NULL,
    address_line2 VARCHAR(200),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postcode VARCHAR(20) NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(200),
    is_default BOOLEAN NOT NULL DEFAULT false,
    instructions TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE customer_default_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    setting_key VARCHAR(100) NOT NULL,
    setting_value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(customer_id, setting_key)
);

CREATE TABLE customer_notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    note TEXT NOT NULL,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_customers_country ON customers(country_code);
CREATE INDEX idx_customers_franchise ON customers(franchise_id);
CREATE INDEX idx_customers_account ON customers(account_number);
CREATE INDEX idx_customer_contacts_customer ON customer_contacts(customer_id);
CREATE INDEX idx_customer_addresses_customer ON customer_addresses(customer_id);
CREATE INDEX idx_customer_default_settings_customer ON customer_default_settings(customer_id);
CREATE INDEX idx_customer_notes_customer ON customer_notes(customer_id);
