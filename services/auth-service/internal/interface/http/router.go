package http

import (
	"net/http"

	"github.com/smexpress/pkg/auth"
	"github.com/smexpress/services/auth-service/internal/interface/http/handler"
)

type Handlers struct {
	Auth       *handler.AuthHandler
	User       *handler.UserHandler
	Role       *handler.RoleHandler
	Permission *handler.PermissionHandler
}

func NewRouter(h Handlers, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()

	authMW := auth.Middleware(jwtSecret)
	adminMW := auth.RequireRole("admin")

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"auth-service"}`))
	})

	// Public routes (no auth)
	mux.HandleFunc("POST /login", h.Auth.Login)
	mux.HandleFunc("POST /refresh", h.Auth.Refresh)

	// Authenticated routes
	mux.Handle("POST /logout", authMW(http.HandlerFunc(h.Auth.Logout)))
	mux.Handle("GET /me", authMW(http.HandlerFunc(h.Auth.GetMe)))
	mux.Handle("POST /change-password", authMW(http.HandlerFunc(h.Auth.ChangePassword)))

	// Admin routes
	mux.Handle("GET /users", authMW(adminMW(http.HandlerFunc(h.User.List))))
	mux.Handle("POST /users", authMW(adminMW(http.HandlerFunc(h.User.Create))))
	mux.Handle("GET /users/{id}", authMW(adminMW(http.HandlerFunc(h.User.Get))))
	mux.Handle("PUT /users/{id}", authMW(adminMW(http.HandlerFunc(h.User.Update))))
	mux.Handle("DELETE /users/{id}", authMW(adminMW(http.HandlerFunc(h.User.Delete))))

	mux.Handle("GET /roles", authMW(http.HandlerFunc(h.Role.List)))
	mux.Handle("POST /roles", authMW(adminMW(http.HandlerFunc(h.Role.Create))))
	mux.Handle("GET /roles/{id}", authMW(http.HandlerFunc(h.Role.Get)))
	mux.Handle("PUT /roles/{id}", authMW(adminMW(http.HandlerFunc(h.Role.Update))))
	mux.Handle("DELETE /roles/{id}", authMW(adminMW(http.HandlerFunc(h.Role.Delete))))
	mux.Handle("PUT /roles/{id}/permissions", authMW(adminMW(http.HandlerFunc(h.Role.SetPermissions))))

	mux.Handle("GET /permissions", authMW(http.HandlerFunc(h.Permission.List)))

	return mux
}
