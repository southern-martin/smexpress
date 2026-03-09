package tenant

import (
	"context"
	"net/http"
)

// Tenant represents the current request's tenant context.
type Tenant struct {
	CountryCode string
	FranchiseID string
	UserID      string
}

type tenantKey struct{}

// WithTenant stores tenant info in context.
func WithTenant(ctx context.Context, t Tenant) context.Context {
	return context.WithValue(ctx, tenantKey{}, t)
}

// FromContext retrieves tenant from context.
func FromContext(ctx context.Context) (Tenant, bool) {
	t, ok := ctx.Value(tenantKey{}).(Tenant)
	return t, ok
}

// MustFromContext retrieves tenant or panics.
func MustFromContext(ctx context.Context) Tenant {
	t, ok := FromContext(ctx)
	if !ok {
		panic("tenant not found in context")
	}
	return t
}

// Middleware extracts tenant information from request headers.
// Expected headers: X-Tenant-Id (country code), X-Franchise-Id (optional).
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			countryCode := r.Header.Get("X-Tenant-Id")
			if countryCode == "" {
				http.Error(w, `{"error":"missing X-Tenant-Id header"}`, http.StatusBadRequest)
				return
			}

			t := Tenant{
				CountryCode: countryCode,
				FranchiseID: r.Header.Get("X-Franchise-Id"),
				UserID:      r.Header.Get("X-User-Id"),
			}

			ctx := WithTenant(r.Context(), t)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
