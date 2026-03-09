---
paths:
  - "services/**/*.go"
  - "pkg/**/*.go"
---
# Go Patterns

## New Service Scaffold

Every service follows this structure:
```
services/<name>/
├── cmd/server/main.go        # entry point: config, DB pool, repos, usecases, handlers, mux, server
└── internal/
    ├── domain/
    │   ├── entity.go          # domain types (structs, enums)
    │   ├── repository.go      # interface for persistence
    │   └── errors.go          # sentinel errors: var ErrNotFound = errors.New("not found")
    ├── usecase/
    │   └── <name>.go          # business logic, depends on repository interface
    ├── interface/http/
    │   ├── handler.go         # HTTP handlers, depends on usecase
    │   ├── dto.go             # request/response DTOs, mapping functions
    │   └── router.go          # func NewRouter(...) http.Handler — registers routes
    └── infrastructure/
        ├── persistence/postgres/
        │   ├── repository.go  # implements domain.Repository using pgxpool
        │   ├── migrations/    # embed.FS SQL migrations (001_init.up.sql, etc.)
        │   └── migrate.go     # RunMigrations() function
        └── config/
            └── config.go      # Load() from env vars
```

## Repository Implementation

- Accept `*pgxpool.Pool` in constructor
- Use `pgx.CollectOneRow` / `pgx.CollectRows` with `pgx.RowToStructByName`
- Always include `country_code` in WHERE clauses for multi-tenancy
- Use `$1, $2, ...` placeholders (not `?`)
- Wrap scan errors: `if errors.Is(err, pgx.ErrNoRows) { return ..., domain.ErrNotFound }`

## Handler Pattern

```go
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")  // Go 1.22+ path params
    entity, err := h.usecase.GetByID(r.Context(), id)
    if err != nil {
        httputil.Error(w, err)  // maps domain errors to HTTP status
        return
    }
    httputil.JSON(w, http.StatusOK, toDTO(entity))
}
```

## Error Wrapping

Always wrap with domain sentinel errors for proper HTTP status mapping:
```go
fmt.Errorf("%w: user %s", domain.ErrNotFound, id)
fmt.Errorf("%w: email already exists", domain.ErrConflict)
```

## Middleware Order

```go
handler := auth.Middleware(cfg.JWTSecret)(router)
handler = logging.HTTPMiddleware(logger)(handler)
```
Logging wraps auth so that all requests (including auth failures) are logged.
