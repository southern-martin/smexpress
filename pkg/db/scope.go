package db

import (
	"context"
	"fmt"
	"strings"
)

// TenantScope holds multi-tenant scoping fields extracted from request context.
type TenantScope struct {
	CountryCode string
	FranchiseID string // UUID string, may be empty
}

type tenantScopeKey struct{}

// WithTenantScope adds tenant scope to a context.
func WithTenantScope(ctx context.Context, scope TenantScope) context.Context {
	return context.WithValue(ctx, tenantScopeKey{}, scope)
}

// GetTenantScope retrieves tenant scope from context.
func GetTenantScope(ctx context.Context) (TenantScope, bool) {
	scope, ok := ctx.Value(tenantScopeKey{}).(TenantScope)
	return scope, ok
}

// ScopeBuilder helps build tenant-scoped SQL WHERE clauses.
type ScopeBuilder struct {
	conditions []string
	args       []any
	argIndex   int
}

// NewScopeBuilder creates a new ScopeBuilder starting at the given argument index.
func NewScopeBuilder(startArgIndex int) *ScopeBuilder {
	return &ScopeBuilder{argIndex: startArgIndex}
}

// ApplyTenantScope adds country_code and optional franchise_id conditions.
func (sb *ScopeBuilder) ApplyTenantScope(scope TenantScope) *ScopeBuilder {
	if scope.CountryCode != "" {
		sb.argIndex++
		sb.conditions = append(sb.conditions, fmt.Sprintf("country_code = $%d", sb.argIndex))
		sb.args = append(sb.args, scope.CountryCode)
	}
	if scope.FranchiseID != "" {
		sb.argIndex++
		sb.conditions = append(sb.conditions, fmt.Sprintf("franchise_id = $%d", sb.argIndex))
		sb.args = append(sb.args, scope.FranchiseID)
	}
	return sb
}

// WhereClause returns the WHERE clause string (without "WHERE" prefix).
func (sb *ScopeBuilder) WhereClause() string {
	if len(sb.conditions) == 0 {
		return "1=1"
	}
	return strings.Join(sb.conditions, " AND ")
}

// Args returns the collected arguments.
func (sb *ScopeBuilder) Args() []any {
	return sb.args
}

// NextArgIndex returns the next available argument index.
func (sb *ScopeBuilder) NextArgIndex() int {
	return sb.argIndex + 1
}
