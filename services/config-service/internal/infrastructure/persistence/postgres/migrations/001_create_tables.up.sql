SET search_path TO imcs_config;

CREATE TABLE system_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    config_key VARCHAR(100) NOT NULL,
    config_value TEXT NOT NULL,
    description TEXT,
    data_type VARCHAR(20) NOT NULL DEFAULT 'string',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, config_key)
);

CREATE TABLE country_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL UNIQUE,
    country_name VARCHAR(100) NOT NULL,
    currency_code VARCHAR(3) NOT NULL,
    currency_symbol VARCHAR(5) NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    date_format VARCHAR(20) NOT NULL DEFAULT 'DD/MM/YYYY',
    weight_unit VARCHAR(5) NOT NULL DEFAULT 'kg',
    dimension_unit VARCHAR(5) NOT NULL DEFAULT 'cm',
    locale VARCHAR(10) NOT NULL DEFAULT 'en',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE feature_flags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    flag_key VARCHAR(100) NOT NULL,
    enabled BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, flag_key)
);

CREATE TABLE holidays (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    holiday_date DATE NOT NULL,
    name VARCHAR(200) NOT NULL,
    is_recurring BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, holiday_date, name)
);

CREATE TABLE sequences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    sequence_type VARCHAR(50) NOT NULL,
    prefix VARCHAR(10),
    current_value BIGINT NOT NULL DEFAULT 0,
    format_pattern VARCHAR(50) NOT NULL DEFAULT '{prefix}{value}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, sequence_type)
);

CREATE INDEX idx_system_configs_country ON system_configs(country_code);
CREATE INDEX idx_feature_flags_country ON feature_flags(country_code);
CREATE INDEX idx_holidays_country_date ON holidays(country_code, holiday_date);
CREATE INDEX idx_sequences_country_type ON sequences(country_code, sequence_type);
