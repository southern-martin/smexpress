package http

import (
	"net/http"

	"github.com/smexpress/services/address-service/internal/interface/http/handler"
)

type Handlers struct {
	Postcode *handler.PostcodeHandler
	Zone     *handler.ZoneHandler
	Region   *handler.RegionHandler
}

func NewRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"address-service"}`))
	})

	// Postcodes
	mux.HandleFunc("GET /postcodes/search", h.Postcode.Search)
	mux.HandleFunc("GET /postcodes/lookup", h.Postcode.Lookup)
	mux.HandleFunc("POST /postcodes", h.Postcode.Create)

	// Zones
	mux.HandleFunc("GET /zones", h.Zone.List)
	mux.HandleFunc("POST /zones", h.Zone.Create)
	mux.HandleFunc("GET /zones/{id}", h.Zone.Get)
	mux.HandleFunc("PUT /zones/{id}", h.Zone.Update)
	mux.HandleFunc("DELETE /zones/{id}", h.Zone.Delete)
	mux.HandleFunc("GET /zones/find", h.Zone.FindZone)
	mux.HandleFunc("PUT /zones/{id}/postcodes", h.Zone.SetPostcodes)

	// Regions
	mux.HandleFunc("GET /regions", h.Region.List)
	mux.HandleFunc("POST /regions", h.Region.Create)
	mux.HandleFunc("GET /regions/{id}", h.Region.Get)
	mux.HandleFunc("PUT /regions/{id}", h.Region.Update)
	mux.HandleFunc("DELETE /regions/{id}", h.Region.Delete)

	return mux
}
