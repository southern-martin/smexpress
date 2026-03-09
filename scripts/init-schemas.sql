-- SmExpress Database Initialization
-- Creates all 17 service schemas and common extensions

-- Extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Service Schemas
CREATE SCHEMA IF NOT EXISTS imcs_auth;
CREATE SCHEMA IF NOT EXISTS imcs_users;
CREATE SCHEMA IF NOT EXISTS imcs_customers;
CREATE SCHEMA IF NOT EXISTS imcs_shipments;
CREATE SCHEMA IF NOT EXISTS imcs_ratings;
CREATE SCHEMA IF NOT EXISTS imcs_carriers;
CREATE SCHEMA IF NOT EXISTS imcs_ecommerce;
CREATE SCHEMA IF NOT EXISTS imcs_invoices;
CREATE SCHEMA IF NOT EXISTS imcs_franchises;
CREATE SCHEMA IF NOT EXISTS imcs_notifications;
CREATE SCHEMA IF NOT EXISTS imcs_reports;
CREATE SCHEMA IF NOT EXISTS imcs_addresses;
CREATE SCHEMA IF NOT EXISTS imcs_documents;
CREATE SCHEMA IF NOT EXISTS imcs_config;
CREATE SCHEMA IF NOT EXISTS imcs_payments;
CREATE SCHEMA IF NOT EXISTS imcs_freight;
CREATE SCHEMA IF NOT EXISTS imcs_live_rating;

-- Grant usage to application user (same as owner in dev)
DO $$
DECLARE
    schema_name TEXT;
BEGIN
    FOR schema_name IN
        SELECT unnest(ARRAY[
            'imcs_auth', 'imcs_users', 'imcs_customers', 'imcs_shipments',
            'imcs_ratings', 'imcs_carriers', 'imcs_ecommerce', 'imcs_invoices',
            'imcs_franchises', 'imcs_notifications', 'imcs_reports', 'imcs_addresses',
            'imcs_documents', 'imcs_config', 'imcs_payments', 'imcs_freight',
            'imcs_live_rating'
        ])
    LOOP
        EXECUTE format('GRANT ALL PRIVILEGES ON SCHEMA %I TO smexpress', schema_name);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON TABLES TO smexpress', schema_name);
        EXECUTE format('ALTER DEFAULT PRIVILEGES IN SCHEMA %I GRANT ALL ON SEQUENCES TO smexpress', schema_name);
    END LOOP;
END $$;
