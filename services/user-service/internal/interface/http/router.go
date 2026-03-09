package http

import (
	"net/http"

	"github.com/smexpress/pkg/auth"
	"github.com/smexpress/services/user-service/internal/interface/http/handler"
)

func NewRouter(h *handler.ProfileHandler, jwtSecret string) *http.ServeMux {
	mux := http.NewServeMux()
	authMW := auth.Middleware(jwtSecret)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"user-service"}`))
	})

	mux.Handle("GET /profile", authMW(http.HandlerFunc(h.GetMyProfile)))
	mux.Handle("PUT /profile", authMW(http.HandlerFunc(h.UpdateMyProfile)))
	mux.Handle("GET /profile/preferences", authMW(http.HandlerFunc(h.GetPreferences)))
	mux.Handle("PUT /profile/preferences/{key}", authMW(http.HandlerFunc(h.SetPreference)))
	mux.Handle("DELETE /profile/preferences/{key}", authMW(http.HandlerFunc(h.DeletePreference)))

	// Admin routes
	mux.Handle("GET /", authMW(http.HandlerFunc(h.List)))
	mux.Handle("GET /{id}/profile", authMW(http.HandlerFunc(h.GetByUserID)))

	return mux
}
