package auth

import (
	"context"
	"net/http"
	"strings"
)

type claimsKey struct{}

// WithClaims stores claims in context.
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// GetClaims retrieves claims from context.
func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	return claims, ok
}

// Middleware returns an HTTP middleware that validates JWT tokens.
func Middleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			claims, err := ParseToken(parts[1], secretKey)
			if err != nil {
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
				return
			}

			ctx := WithClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole returns middleware that checks for a specific role.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaims(r.Context())
			if !ok {
				http.Error(w, `{"error":"no claims in context"}`, http.StatusUnauthorized)
				return
			}

			for _, claimRole := range claims.Roles {
				if claimRole == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, `{"error":"insufficient permissions"}`, http.StatusForbidden)
		})
	}
}
