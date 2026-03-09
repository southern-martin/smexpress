SET search_path TO imcs_franchises;

CREATE TABLE franchises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(20) NOT NULL,
    contact_name VARCHAR(200),
    email VARCHAR(200),
    phone VARCHAR(20),
    address_line1 VARCHAR(200),
    address_line2 VARCHAR(200),
    city VARCHAR(100),
    state VARCHAR(100),
    postcode VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT true,
    commission_rate NUMERIC(5,4) NOT NULL DEFAULT 0,
    parent_franchise_id UUID REFERENCES franchises(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, code)
);

CREATE TABLE franchise_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL REFERENCES franchises(id),
    setting_key VARCHAR(100) NOT NULL,
    setting_value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(franchise_id, setting_key)
);

CREATE TABLE territories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL REFERENCES franchises(id),
    country_code VARCHAR(2) NOT NULL,
    name VARCHAR(200) NOT NULL,
    postcode_from VARCHAR(20),
    postcode_to VARCHAR(20),
    state VARCHAR(100),
    is_exclusive BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE franchise_ledgers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL UNIQUE REFERENCES franchises(id),
    country_code VARCHAR(2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    balance NUMERIC(14,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE franchise_ledger_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ledger_id UUID NOT NULL REFERENCES franchise_ledgers(id),
    entry_type VARCHAR(20) NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    balance_after NUMERIC(14,2) NOT NULL,
    description TEXT,
    reference_type VARCHAR(50),
    reference_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE franchise_withdrawals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL REFERENCES franchises(id),
    country_code VARCHAR(2) NOT NULL,
    amount NUMERIC(14,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    requested_by UUID NOT NULL,
    approved_by UUID,
    bank_account_name VARCHAR(200),
    bank_account_number VARCHAR(50),
    bank_bsb VARCHAR(20),
    notes TEXT,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

CREATE TABLE franchise_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    franchise_id UUID NOT NULL REFERENCES franchises(id),
    action VARCHAR(50) NOT NULL,
    changed_by UUID NOT NULL,
    changes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_franchises_country ON franchises(country_code);
CREATE INDEX idx_franchise_settings_franchise ON franchise_settings(franchise_id);
CREATE INDEX idx_territories_franchise ON territories(franchise_id);
CREATE INDEX idx_territories_country ON territories(country_code);
CREATE INDEX idx_franchise_ledgers_franchise ON franchise_ledgers(franchise_id);
CREATE INDEX idx_franchise_ledger_entries_ledger ON franchise_ledger_entries(ledger_id);
CREATE INDEX idx_franchise_withdrawals_franchise ON franchise_withdrawals(franchise_id);
CREATE INDEX idx_franchise_withdrawals_status ON franchise_withdrawals(status);
CREATE INDEX idx_franchise_history_franchise ON franchise_history(franchise_id);
