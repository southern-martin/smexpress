package http

import (
	"net/http"

	"github.com/smexpress/services/franchise-service/internal/interface/http/handler"
)

type Handlers struct {
	Franchise  *handler.FranchiseHandler
	Territory  *handler.TerritoryHandler
	Ledger     *handler.LedgerHandler
	Withdrawal *handler.WithdrawalHandler
}

func NewRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"franchise-service"}`))
	})

	// Franchises
	mux.HandleFunc("GET /franchises", h.Franchise.List)
	mux.HandleFunc("POST /franchises", h.Franchise.Create)
	mux.HandleFunc("GET /franchises/{id}", h.Franchise.Get)
	mux.HandleFunc("PUT /franchises/{id}", h.Franchise.Update)
	mux.HandleFunc("PATCH /franchises/{id}/activate", h.Franchise.Activate)
	mux.HandleFunc("PATCH /franchises/{id}/deactivate", h.Franchise.Deactivate)

	// Franchise settings
	mux.HandleFunc("GET /franchises/{id}/settings", h.Franchise.GetSettings)
	mux.HandleFunc("PUT /franchises/{id}/settings", h.Franchise.SetSetting)

	// Territories
	mux.HandleFunc("GET /franchises/{id}/territories", h.Territory.ListByFranchise)
	mux.HandleFunc("POST /franchises/{id}/territories", h.Territory.Create)
	mux.HandleFunc("GET /franchises/{id}/territories/{territoryId}", h.Territory.Get)
	mux.HandleFunc("PUT /franchises/{id}/territories/{territoryId}", h.Territory.Update)
	mux.HandleFunc("DELETE /franchises/{id}/territories/{territoryId}", h.Territory.Delete)
	mux.HandleFunc("GET /territories/search", h.Territory.FindByPostcode)

	// Ledger
	mux.HandleFunc("GET /franchises/{id}/ledger", h.Ledger.GetBalance)
	mux.HandleFunc("GET /franchises/{id}/ledger/entries", h.Ledger.ListEntries)
	mux.HandleFunc("POST /franchises/{id}/ledger/credit", h.Ledger.Credit)

	// Withdrawals
	mux.HandleFunc("POST /withdrawals", h.Withdrawal.Request)
	mux.HandleFunc("GET /withdrawals/pending", h.Withdrawal.ListPending)
	mux.HandleFunc("GET /franchises/{id}/withdrawals", h.Withdrawal.ListByFranchise)
	mux.HandleFunc("GET /withdrawals/{withdrawalId}", h.Withdrawal.Get)
	mux.HandleFunc("PATCH /withdrawals/{withdrawalId}/approve", h.Withdrawal.Approve)
	mux.HandleFunc("PATCH /withdrawals/{withdrawalId}/reject", h.Withdrawal.Reject)
	mux.HandleFunc("PATCH /withdrawals/{withdrawalId}/process", h.Withdrawal.Process)

	return mux
}
