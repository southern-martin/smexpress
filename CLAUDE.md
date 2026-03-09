# smexpress

IMCS PHP-to-Go microservices migration. Logistics/shipping platform (42+ carriers, 13 countries).
17 Go microservices + Kong Gateway + React/Tailwind frontend.

## Build & Run

```bash
# Build single service
go build ./services/<name>/...

# Build all services
make build

# Run a service
make run-svc SVC=auth-service

# Tests
make test                        # all tests
make test-svc SVC=auth-service   # single service
go test -v ./pkg/...             # shared packages only

# Lint & format
make vet                         # go vet
make fmt                         # gofmt

# Infrastructure (PostgreSQL, NATS, Kong)
make docker-up
make docker-down
```

### Frontend

```bash
cd frontend && pnpm install      # install deps
cd frontend && pnpm dev          # dev servers (admin :3000, customer :3001)
cd frontend && pnpm build        # production build
```

## Go Workspace

- `go.work` at root, each service has minimal `go.mod` (module name + `go 1.24.0`)
- No `go mod tidy` needed — workspace resolves all deps from `pkg/go.mod`
- Services don't need explicit `require` directives for `pkg/`

## Architecture (per service)

```
services/<name>/
├── cmd/server/main.go
└── internal/
    ├── domain/           # entities, repository interfaces, domain errors
    ├── usecase/          # application logic, orchestrates domain
    ├── interface/http/   # handler.go, dto.go, router.go
    └── infrastructure/
        ├── persistence/postgres/   # repository impls, migrations/
        └── config/                 # env-based config
```

## Code Conventions

### HTTP Layer
- Go 1.22+ `net/http.ServeMux` with method patterns: `"GET /path"`, `"POST /path/{id}"`
- Handler flow: decode request -> call use case -> map entity to DTO -> `httputil.JSON()` / `Created()` / `NoContent()`
- Middleware chain: `logging.HTTPMiddleware(logger)(mux)` then `auth.Middleware(secret)(handler)`

### Database
- pgxpool (pgx/v5), raw SQL — no ORM
- Schema-per-service in shared PostgreSQL (e.g., `imcs_auth`, `imcs_config`)
- Multi-tenancy: `country_code` column on every table, propagated via `pkg/tenant`
- Migrations: `embed.FS` in `infrastructure/persistence/postgres/migrations/`, `RunMigrations()` function
- Pagination: `db.Page{Number, Size}` / `db.PagedResult[T]` / `page.LimitOffsetClause(argStart)`

### Error Handling
- Wrap domain errors: `fmt.Errorf("%w: details", domainerr.ErrNotFound)`
- Duplicate key check: `isDuplicateKey(err)` using `strings.Contains(err.Error(), "duplicate key")`

### Auth
- JWT via `pkg/auth` — `GetClaims(ctx)` returns `(*Claims, bool)`

### Events
- NATS JetStream via `pkg/messaging` — Publisher/Subscriber pattern

## Shared Packages (`pkg/`)

auth, db, httputil, logging, messaging, money, tenant, testutil, proto

## Service Ports

| Service | Port | Schema |
|---------|------|--------|
| auth | 8081 | imcs_auth |
| config | 8082 | imcs_config |
| user | 8083 | imcs_users |
| franchise | 8084 | imcs_franchises |
| customer | 8085 | imcs_customers |
| address | 8086 | imcs_addresses |
| carrier | 8087 | imcs_carriers |
| rating | 8088 | imcs_rating |
| shipment | 8089 | imcs_shipments |
| document | 8090 | imcs_documents |
| notification | 8091 | imcs_notifications |
| ecommerce | 8092 | imcs_ecommerce |
| invoice | 8093 | imcs_invoices |
| reporting | 8094 | imcs_reporting |
| payment | 8095 | imcs_payments |
| freight | 8096 | imcs_freight |
| live-rating | 8097 | imcs_live_rating |

## Frontend

- Monorepo: pnpm + Turborepo at `frontend/`
- Apps: `apps/admin` (port 3000), `apps/customer` (port 3001)
- Shared: @smexpress/ui, @smexpress/auth, @smexpress/api-client, @smexpress/config, @smexpress/i18n
- Stack: React 18.3, Vite, TanStack Query, Zustand, Tailwind CSS, React Router 6

## Git Workflow (Git Flow)

- **Branches**: `main` (production), `develop` (integration), `feature/*`, `release/*`, `hotfix/*`
- **Never push directly** to `main` or `develop`
- **Feature work**: branch from `develop` → `feature/<name>` → merge to `develop` with `--no-ff`
- **Releases**: branch from `develop` → `release/<version>` → merge to both `main` and `develop` with `--no-ff`, tag `main`
- **Hotfixes**: branch from `main` → `hotfix/<name>` → merge to both `main` and `develop` with `--no-ff`, tag `main`
- **Commit format**: `<type>(<scope>): <description>` — types: feat, fix, refactor, test, docs, chore
- **Delete feature branches** after merge

## Key Files

- Migration plan: `plans/migration-plan.md`
- Docker infra: `deployments/docker-compose.yml`
- Kong config: `kong/kong.yml`
- DB schema init: `scripts/init-schemas.sql`
