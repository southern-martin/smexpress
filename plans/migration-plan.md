# IMCS PHP to Go Microservices Migration Plan

## Context

Migrating the IMCS logistics/shipping management platform from a legacy PHP monolith (75+ modules, 729+ DB mappings, 42+ carrier integrations, 12+ ecommerce platforms, 13 countries) to a modern stack. The old PHP system has no database and cannot be run — all business logic must be reverse-engineered from source code.

**Target Stack:** Go (clean architecture) + Kong API Gateway + React/Tailwind frontend
**Target Location:** `/opt/openAi/smexpress/`

---

## Architecture Overview

```
Browser --> Kong Gateway (port 8000)
              |
              |--> 17 Go Microservices (gRPC inter-service)
              |
              |--> NATS JetStream (async events)
              |
PostgreSQL (schema-per-service) + Redis (cache/sessions) + MinIO (documents)
```

**Multi-tenancy:** Single deployment, row-level isolation via `country_code` column (AU, CA, DE, FR, HK, IN, KR, MA, NL, NZ, UK, VN, ZA). Configuration-driven country behavior.

---

## Part 1: Monorepo Structure

```
smexpress/
├── services/                    # 17 Go microservices
│   ├── auth-service/
│   ├── user-service/
│   ├── customer-service/
│   ├── shipment-service/
│   ├── rating-service/
│   ├── carrier-service/
│   ├── ecommerce-service/
│   ├── invoice-service/
│   ├── franchise-service/
│   ├── notification-service/
│   ├── reporting-service/
│   ├── address-service/
│   ├── document-service/
│   ├── config-service/
│   ├── payment-service/
│   ├── freight-service/
│   └── live-rating-service/
├── pkg/                         # Shared Go libraries
│   ├── auth/                    # JWT parsing, middleware
│   ├── httputil/                # JSON response helpers, pagination
│   ├── db/                      # PostgreSQL connection, multi-tenant scoping
│   ├── messaging/               # NATS JetStream pub/sub
│   ├── logging/                 # Structured logging (zerolog/slog)
│   ├── tenant/                  # Tenant context propagation
│   ├── money/                   # Currency-safe money type
│   ├── proto/                   # Shared protobuf definitions
│   └── testutil/                # Test helpers, testcontainers
├── api/
│   ├── proto/                   # All .proto files
│   └── openapi/                 # OpenAPI specs
├── kong/
│   └── kong.yml                 # Declarative Kong config
├── frontend/                    # React + Tailwind (Turborepo monorepo)
│   ├── apps/
│   │   ├── admin/               # Admin portal SPA
│   │   └── customer/            # Customer portal SPA
│   └── packages/
│       ├── ui/                  # Shared component library
│       ├── api-client/          # Axios + API endpoints
│       ├── auth/                # JWT management, ProtectedRoute
│       ├── i18n/                # Internationalization
│       └── config/              # Tenant configs, feature flags
├── scripts/
│   ├── init-schemas.sql
│   └── generate-proto.sh
├── deployments/
│   └── docker-compose.yml
├── go.work
├── Makefile
└── README.md
```

---

## Part 2: 17 Microservices

### Service Decomposition

| # | Service | Responsibility | Key Entities | Dependencies |
|---|---------|---------------|--------------|-------------|
| 1 | **auth-service** | Login, JWT, roles, permissions | users, roles, permissions, refresh_tokens | None |
| 2 | **config-service** | System config, feature flags, lookups | system_config, country_config, feature_flags, holidays, areas, industries, hts_goods, commodity_classifications, insurance_covers, languages, language_values, sequences, boards, board_widgets | None |
| 3 | **user-service** | User profiles, preferences | user_profiles, user_preferences | auth |
| 4 | **franchise-service** | Franchise/branch management | franchises, franchise_settings, territories, franchise_ledgers, franchise_ledger_entries, franchise_withdrawals, franchise_history, franchise_permissions | auth |
| 5 | **customer-service** | Customer accounts, contacts, credit | customers, customer_contacts, customer_addresses, customer_default_settings, customer_commodity_skus, customer_importers, customer_payment_reserved, customer_notes, customer_automated_rules, automated_rule_criteria, automated_rule_actions, customer_carrier_accounts (BYOA) | auth, franchise |
| 6 | **address-service** | Postcode lookup, validation, zones | postcodes, zones, regions, eu_countries, benelux_countries, address_validations | None |
| 7 | **rating-service** | Rate cards, quoting, surcharges | rate_cards, rate_zones, surcharge_rules, quotes, quote_logs, base_rates, base_rate_details, base_rate_templates, franchise_base_rates, rate_sheets, rate_sheet_details, rate_sheet_rows, rate_sheet_columns, transit_times, transit_time_mappings, accessorials, accessorial_templates, accessorial_template_details, franchise_accessorial_rates, customer_accessorial_rates | carrier, customer |
| 8 | **carrier-service** | Carrier API integration layer | carriers, carrier_services, carrier_credentials, carrier_service_mappings, customer_service_mappings, franchise_carriers | External APIs |
| 9 | **shipment-service** | Shipment lifecycle, tracking, manifests | shipments, shipment_items, shipment_events, shipment_billing, shipment_billing_charges, shipment_accessorial_charges, shipment_history, shipment_addresses, shipment_extra_info, pieces, piece_tracking, batch_shipments, batch_shipment_items, batch_shipment_pieces, commercial_invoices, scheduled_collections, reschedule_collections, manifests | carrier, rating, address, document |
| 10 | **document-service** | PDF generation, label storage | document_templates, documents | MinIO/S3 |
| 11 | **notification-service** | Email, SMS, webhooks | notification_templates, notification_logs, email_queue, quote_log_emails, adjustment_notifications, reminder_letter_logs | None (event-driven) |
| 12 | **ecommerce-service** | Platform integrations, order import | ecommerce_connections, ecommerce_orders, ecommerce_order_items, ecommerce_account_settings (supports: Shopify, WooCommerce, Magento, Magento2, Amazon, eBay, Etsy, BigCommerce, Ecwid, Wix, Squarespace, PrestaShop) | shipment |
| 13 | **invoice-service** | Billing, payments, credit notes | invoices, invoice_lines, payments, credit_notes, credit_note_details, statement_invoices, invoice_freight_credits, adjustments, adjustment_carriers, adjustment_settings, adjustment_notifications, adjustment_transaction_logs, credit_card_fees, overpayments, transactions, transaction_details, reminder_letters | shipment, customer |
| 14 | **reporting-service** | Dashboards, reports, aggregation | reports, dashboard_widgets, aggregated_stats | Event-driven reads |
| 15 | **payment-service** | Payment processing, Adyen/gateway integrations | payment_gateways, payment_gateway_accounts, payment_transactions, payment_callbacks | accounting |
| 16 | **freight-service** | Freight shipping, freight pricing, freight billing | freight_shipments, freight_pricing, freight_billing, freight_transactions | carrier, shipment |
| 17 | **live-rating-service** | Real-time rate widget for ecommerce stores | live_rating_stores, live_rating_boxes, live_rating_accounts, live_rating_products | rating, carrier |

### Clean Architecture per Service

```
service-name/
├── cmd/server/main.go              # Entry point, dependency injection
├── internal/
│   ├── domain/                     # DOMAIN (no external imports)
│   │   ├── entity/                 # Business entities
│   │   ├── valueobject/            # Value objects (Address, Money, Weight)
│   │   ├── repository/            # Repository interfaces (ports)
│   │   ├── event/                  # Domain events
│   │   └── errors/                 # Domain errors
│   ├── usecase/                    # USE CASES (depends on domain only)
│   │   └── shipment/
│   │       ├── create_shipment.go  # One file per use case
│   │       ├── book_shipment.go
│   │       └── list_shipments.go
│   ├── interface/                  # INTERFACE LAYER
│   │   ├── http/
│   │   │   ├── handler/           # HTTP handlers
│   │   │   ├── middleware/        # Auth, logging, tenant
│   │   │   ├── dto/               # Request/response DTOs
│   │   │   └── router.go
│   │   ├── grpc/                  # gRPC handlers (inter-service)
│   │   └── consumer/             # NATS event consumers
│   └── infrastructure/            # INFRASTRUCTURE
│       ├── persistence/postgres/  # Repository implementations + migrations
│       ├── persistence/redis/     # Cache implementations
│       ├── external/              # gRPC clients to other services
│       ├── messaging/             # NATS publisher/subscriber
│       └── config/                # Config loading
├── api/proto/                     # Service protobuf definitions
├── Dockerfile
├── go.mod
└── Makefile
```

**Layer rules:** domain -> no imports | usecase -> domain only | interface -> usecase+domain | infrastructure -> implements domain interfaces

### Inter-Service Communication

**Sync (gRPC):**
- shipment->carrier (book)
- shipment->rating (quote)
- shipment->address (validate)
- shipment->document (commercial invoice PDF)
- rating->customer (custom rates)
- rating->carrier (carrier service availability)
- live-rating->rating (real-time quotes)
- live-rating->carrier (service check)
- invoice->franchise (franchise ledger updates)

**Async (NATS JetStream):**
- `shipment.booked` -> notification, ecommerce, reporting
- `shipment.created` -> reporting, notification
- `batch_shipment.processed` -> notification, reporting
- `tracking.updated` -> shipment, notification
- `invoice.generated` -> notification
- `adjustment.created` -> invoice, notification
- `payment.received` -> invoice, notification
- `franchise_withdrawal.requested` -> invoice, notification
- `ecommerce.order_imported` -> shipment, notification

### Key API Endpoints (via Kong)

**Auth:** `POST /auth/login`, `POST /auth/refresh`, `GET /auth/me`
**Shipments:** `GET/POST /shipments`, `POST /shipments/{id}/book`, `GET /shipments/{id}/tracking`
**Batch Shipments:** `POST /shipments/batch`, `GET /shipments/batch/{id}/status`
**Quotes:** `POST /quotes`, `POST /quotes/compare`, `GET /quotes/logs`, `POST /quotes/email`
**Customers:** `GET/POST /customers`, `GET /customers/{id}/addresses`
**Automated Rules:** `GET/POST /customers/{id}/rules`
**Invoices:** `GET/POST /invoices`, `POST /billing-runs`
**Adjustments:** `GET/POST /adjustments`, `POST /adjustments/{id}/apply`
**Carriers:** `GET /carriers`, `POST /carriers/{code}/book` (internal)
**Ecommerce:** `POST /ecommerce/connections`, `POST /ecommerce/orders/bulk-import`
**Rate Sheets:** `GET/POST /rate-sheets`, `POST /rate-sheets/import`, `GET /rate-sheets/export`
**Franchise:** `GET /franchises/{id}/ledger`, `POST /franchises/{id}/withdrawals`
**Live Rating:** `POST /live-rating/quote`, `GET /live-rating/stores`, `POST /live-rating/stores`
**Payments:** `POST /payments/process`, `POST /payments/callback/{provider}`
**Freight:** `POST /freight/shipments`, `GET /freight/pricing`
**Files:** `POST /files/upload`, `GET /files/{id}`

---

## Part 3: Database Design

**Strategy:** Schema-per-service in shared PostgreSQL cluster (can split to separate DBs later without code changes).

17 schemas: `imcs_auth`, `imcs_users`, `imcs_customers`, `imcs_shipments`, `imcs_ratings`, `imcs_carriers`, `imcs_ecommerce`, `imcs_invoices`, `imcs_franchises`, `imcs_notifications`, `imcs_reports`, `imcs_addresses`, `imcs_documents`, `imcs_config`, `imcs_payments`, `imcs_freight`, `imcs_live_rating`

**Total tables:** ~450+

**Multi-tenant pattern:** Every table includes `country_code VARCHAR(2) NOT NULL` + optional `franchise_id UUID`. Repository layer auto-scopes queries from tenant context.

**Redis:** Session/token storage (TTL 24h) + application cache (rates, config, tracking — TTL varies).

---

## Part 4: Kong Gateway

- Declarative (DB-less) mode via `kong.yml`
- JWT plugin on all routes except `/auth/login`, `/auth/password/reset`, ecommerce webhooks
- Unprotected endpoints:
  - Ecommerce platform webhooks (`/webhooks/shopify`, `/webhooks/woocommerce`, etc.)
  - Payment gateway callbacks (`/payments/callback`)
  - Live rating widget API (`/live-rating/quote`)
- Rate limiting: 60/min global, 10/min on auth endpoints
- CORS configured for frontend domains
- Correlation ID propagation (`X-Request-Id`)
- Request size limit: 10MB

---

## Part 5: Frontend Architecture

### Tech Stack

| Concern | Choice |
|---------|--------|
| Monorepo | Turborepo + pnpm |
| Framework | React 18+ |
| Routing | React Router v6 (lazy-loaded) |
| Styling | Tailwind CSS |
| Server state | TanStack Query v5 |
| Client state | Zustand |
| Forms | React Hook Form + Zod |
| Tables | TanStack Table v8 |
| Charts | Recharts |
| i18n | react-i18next |
| HTTP client | Axios (interceptors for JWT auto-refresh) |
| Build | Vite + SWC |
| Deploy | Docker + nginx |

### Two SPAs

**Admin Portal** (`apps/admin`): Dashboard, rate management, **rate sheet import/export**, user/role management, franchise management, **franchise ledger/withdrawals**, invoice management, credit notes, **adjustments (8 adjustment types)**, **payment/overpayment management**, system config, **holiday management**, **area/territory management**, **reminder letter management**, **quote log management**, **freight management**, **carrier configuration per country**, **language management**, **live rating store management**

**Customer Portal** (`apps/customer`): Shipment creation/tracking, **batch shipment creation/upload**, quote requests, invoice viewing, address book, ecommerce order management, **automated rules management**, **commercial invoice management**, **scheduled collection management**, **BYOA carrier account setup**, account settings, **live rating widget configuration**

### Auth Flow
- JWT access token (in-memory) + refresh token (localStorage)
- Axios interceptor handles auto-refresh on 401
- `ProtectedRoute` + `RoleGuard` components for route protection
- `X-Tenant-ID` header sent with every request

### Multi-tenancy
- Tenant config map defines per-country: locale, currency, date format, weight/dimension units, feature flags
- Country switcher in top bar, persisted in Zustand
- Feature flags control UI sections per country (e.g., franchise features only for AU/UK)

---

## Part 6: Infrastructure (docker-compose)

Services in docker-compose for local dev:
- **PostgreSQL 16** (port 5432)
- **Redis 7** (port 6379)
- **NATS JetStream** (port 4222)
- **Kong 3.6** (port 8000 proxy, 8001 admin)
- **MinIO** (port 9000 — S3-compatible document storage)
- **17 Go services** (ports 8081-8097)
- **Frontend** (port 3000)

---

## Part 7: Phased Implementation Roadmap

### Phase 0: Foundation (Weeks 1-3)
- Initialize Go workspace monorepo (`go.work`)
- Build shared packages: `pkg/auth`, `pkg/httputil`, `pkg/db`, `pkg/messaging`, `pkg/logging`, `pkg/tenant`
- Create `docker-compose.yml` with PostgreSQL, Redis, NATS, Kong, MinIO
- Create `init-schemas.sql` for all 17 schemas
- Build service template/skeleton
- Set up CI pipeline (lint, test, build)
- Configure Kong declarative config
- Initialize frontend monorepo (Turborepo, pnpm, shared packages)
- Build `packages/api-client`, `packages/auth`, `packages/config`, `packages/ui` scaffolds

### Phase 1: Auth + Config + Users (Weeks 4-7)
- **config-service** — country configs, feature flags, lookups
- **auth-service** — login, JWT, roles/permissions CRUD
- **user-service** — user profiles, preferences
- **franchise-service** — franchise CRUD
- **Frontend:** Login page, dashboard shell, admin user/franchise/config screens

### Phase 2: Customers + Addresses (Weeks 8-10)
- **customer-service** — CRUD, contacts, addresses, credit
- **address-service** — postcode lookup, validation, zone resolution
- **Frontend:** Customer list/detail, creation wizard, address lookup components

### Phase 3: Rating + Carrier Core (Weeks 11-16)
- **rating-service** — rate cards, zones, surcharges, quoting engine
- **carrier-service** — adapter framework + first 5 carriers: DHL (International+Domestic), DPD, UPS, FedEx, TNT
- **Frontend:** Quote comparison tool, rate card management, carrier config screens

### Phase 4: Shipments + Documents (Weeks 17-22) **[MVP]**
- **shipment-service** — full lifecycle (create, book, track, manifest), batch shipments, commercial invoices, scheduled collections, adjustment system
- **document-service** — label PDF generation, customs docs
- **Frontend:** Shipment wizard, list/detail, label print, manifest, bulk upload, tracking

### Phase 5: Notifications + Ecommerce (Weeks 23-28)
- **notification-service** — email templates, sending, webhook delivery
- **ecommerce-service** — Shopify, WooCommerce, Magento, Amazon adapters
- **Frontend:** Notification templates, ecommerce connection wizard, order import

### Phase 6: Invoicing + Reporting + Franchise Enhancements (Weeks 29-34)
- **invoice-service** — billing runs, invoices, payments, credit notes, **adjustments, statement invoices, overpayment management, reminder letters**
- **reporting-service** — dashboards, reports, data aggregation, **franchise ledger reports, carrier payment reports**
- **franchise-service enhancements** — franchise ledger, withdrawals, commission tracking
- **Frontend:** Invoice management, billing runs, dashboard widgets, report builder

### Phase 7: Remaining Carriers + Ecommerce (Weeks 35-44)
- **UK carriers:** UKMail, Whistl, Yodel, Hermes, Royal Mail
- **AU carriers:** StarTrack, Toll Priority, Toll IPEC, Direct Courier, Northline
- **NZ carriers:** Direct Courier
- **CA carriers:** Canpar, Purolator, Loomis
- **FR carriers:** Chronopost, GLS, DPD3
- **NL carriers:** DHL Parcel, DPD2
- **DE carriers:** DPD
- **IN carriers:** Blue Dart
- **ZA carriers:** JKJ, RAM
- **Additional:** USPS, DHL eCommerce, FedEx BYOA variants
- **Remaining ecommerce:** eBay, BigCommerce, Etsy, Ecwid, Wix, Squarespace, PrestaShop, Magento2

### Phase 7.5: Specialized Features (Weeks 45-50)
- **payment-service** — Adyen, payment gateway framework
- **freight-service** — freight pricing, freight billing, freight transactions
- **live-rating-service** — store integration, real-time widget API
- **BYOA module** — customer carrier account management
- **Insurance integration** — FreightCover third-party insurance
- **Zonos integration** — duties/taxes calculation
- **Automated rules engine** — criteria/action framework
- **E-Invoice/E-Credit** — India GST compliance, Vietnam tax compliance
- **Frontend:** Payment management, freight management, live rating admin, BYOA setup, automated rules UI

### Phase 9: Production Hardening (Weeks 55-62)
- Load testing, OpenTelemetry tracing, Prometheus/Grafana
- Security audit, Kubernetes manifests, CI/CD pipeline
- Documentation, runbooks

---

## Part 8: Verification Plan

### Per-Service Testing
1. **Unit tests:** Domain entities, use cases (mock repositories)
2. **Integration tests:** Repository implementations against test PostgreSQL (testcontainers-go)
3. **API tests:** HTTP handler tests with httptest
4. **Contract tests:** gRPC proto compatibility checks

### End-to-End Testing
1. `docker-compose up` — all services, infra, frontend start
2. Login via frontend -> verify JWT flow through Kong
3. Create customer -> create quote -> book shipment -> verify tracking
4. Import ecommerce order -> convert to shipment
5. Run billing -> verify invoice generation
6. Verify multi-tenant: switch country, confirm correct carriers/rates/config
7. Batch shipment: Upload CSV -> Process batch -> Verify all shipments created
8. Adjustment: Create adjustment -> Apply to invoice -> Verify ledger
9. Franchise ledger: Book shipment -> Verify franchise commission -> Process withdrawal
10. Live rating: Configure store -> Request quote via widget -> Verify rates
11. BYOA: Configure customer carrier account -> Book shipment with customer's account
12. Automated rules: Create rule -> Import ecommerce order -> Verify rule applied
13. Multi-country: Test carriers across AU, UK, CA, FR, NL, DE, NZ, IN, ZA

### Frontend Testing
1. **Vitest + Testing Library:** Component and hook unit tests
2. **Playwright:** E2E tests for critical flows (login, shipment creation, quote)

---

## Summary

| Metric | Value |
|--------|-------|
| Total microservices | 17 |
| Database schemas | 17 |
| Estimated tables | ~450+ |
| Carrier integrations | 42+ |
| Ecommerce platforms | 12+ |
| Countries supported | 13 (AU, CA, DE, FR, HK, IN, KR, MA, NL, NZ, UK, VN, ZA) |
| Implementation phases | 9 (including Phase 7.5) |
| Estimated timeline | ~62 weeks |
