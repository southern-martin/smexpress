SET search_path TO imcs_addresses;

CREATE TABLE postcodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    postcode VARCHAR(20) NOT NULL,
    suburb VARCHAR(200),
    city VARCHAR(200),
    state VARCHAR(100),
    state_code VARCHAR(10),
    latitude NUMERIC(10,7),
    longitude NUMERIC(10,7),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    zone_name VARCHAR(100) NOT NULL,
    zone_code VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, zone_code)
);

CREATE TABLE zone_postcodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    zone_id UUID NOT NULL REFERENCES zones(id) ON DELETE CASCADE,
    postcode_from VARCHAR(20) NOT NULL,
    postcode_to VARCHAR(20) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE regions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(2) NOT NULL,
    name VARCHAR(200) NOT NULL,
    code VARCHAR(20) NOT NULL,
    parent_region_id UUID REFERENCES regions(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(country_code, code)
);

CREATE INDEX idx_postcodes_country_postcode ON postcodes(country_code, postcode);
CREATE INDEX idx_postcodes_country_suburb ON postcodes(country_code, suburb);
CREATE INDEX idx_zones_country ON zones(country_code);
CREATE INDEX idx_zone_postcodes_zone ON zone_postcodes(zone_id);
CREATE INDEX idx_regions_country ON regions(country_code);
