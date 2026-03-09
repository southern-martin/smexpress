package http

import (
	"net/http"

	"github.com/smexpress/services/config-service/internal/interface/http/handler"
)

type Handlers struct {
	SystemConfig *handler.SystemConfigHandler
	CountryConfig *handler.CountryConfigHandler
	FeatureFlag  *handler.FeatureFlagHandler
	Holiday      *handler.HolidayHandler
	Sequence     *handler.SequenceHandler
}

func NewRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"config-service"}`))
	})

	// System configs
	mux.HandleFunc("GET /system-configs", h.SystemConfig.List)
	mux.HandleFunc("POST /system-configs", h.SystemConfig.Create)
	mux.HandleFunc("GET /system-configs/{id}", h.SystemConfig.Get)
	mux.HandleFunc("PUT /system-configs/{id}", h.SystemConfig.Update)
	mux.HandleFunc("DELETE /system-configs/{id}", h.SystemConfig.Delete)

	// Country configs
	mux.HandleFunc("GET /countries", h.CountryConfig.List)
	mux.HandleFunc("POST /countries", h.CountryConfig.Create)
	mux.HandleFunc("GET /countries/{code}", h.CountryConfig.Get)
	mux.HandleFunc("PUT /countries/{code}", h.CountryConfig.Update)

	// Feature flags
	mux.HandleFunc("GET /feature-flags", h.FeatureFlag.List)
	mux.HandleFunc("POST /feature-flags", h.FeatureFlag.Create)
	mux.HandleFunc("PUT /feature-flags/{id}", h.FeatureFlag.Update)
	mux.HandleFunc("PATCH /feature-flags/{id}/toggle", h.FeatureFlag.Toggle)

	// Holidays
	mux.HandleFunc("GET /holidays", h.Holiday.List)
	mux.HandleFunc("POST /holidays", h.Holiday.Create)
	mux.HandleFunc("DELETE /holidays/{id}", h.Holiday.Delete)

	// Sequences
	mux.HandleFunc("GET /sequences", h.Sequence.List)
	mux.HandleFunc("POST /sequences", h.Sequence.Create)
	mux.HandleFunc("POST /sequences/{type}/next", h.Sequence.NextValue)

	return mux
}
