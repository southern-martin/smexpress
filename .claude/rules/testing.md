---
paths:
  - "**/*_test.go"
---
# Testing Conventions

## Commands

```bash
make test                        # all tests
make test-svc SVC=auth-service   # single service
go test -v -run TestFuncName ./services/auth-service/...  # single test
```

## Guidelines

- Test files live next to the code they test (`repository_test.go` beside `repository.go`)
- Use table-driven tests for multiple cases
- Use `testutil` package from `pkg/testutil` for shared test helpers
- Mock at the repository interface level — usecases accept interfaces
- Name tests: `Test<Function>_<scenario>` (e.g., `TestGetByID_NotFound`)
- Prefer `t.Run(name, func(t *testing.T) {...})` subtests for readability

## Integration Tests

- Use build tag `//go:build integration` for tests needing a real database
- Use `testutil.SetupTestDB()` for database setup/teardown
- Run with: `go test -tags=integration ./services/<name>/...`
