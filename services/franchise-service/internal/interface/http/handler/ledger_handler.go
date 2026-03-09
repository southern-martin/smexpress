package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
	"github.com/smexpress/services/franchise-service/internal/interface/http/dto"
	"github.com/smexpress/services/franchise-service/internal/usecase"
)

type LedgerHandler struct {
	uc *usecase.LedgerUseCase
}

func NewLedgerHandler(uc *usecase.LedgerUseCase) *LedgerHandler {
	return &LedgerHandler{uc: uc}
}

func (h *LedgerHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	ledger, err := h.uc.GetBalance(r.Context(), franchiseID)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "ledger not found")
			return
		}
		httputil.InternalError(w, "failed to get balance")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.LedgerFromEntity(ledger))
}

func (h *LedgerHandler) ListEntries(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.ListEntries(r.Context(), franchiseID, page)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "ledger not found")
			return
		}
		httputil.InternalError(w, "failed to list entries")
		return
	}

	type pagedResponse struct {
		Items      []dto.LedgerEntryResponse `json:"items"`
		TotalCount int64                     `json:"total_count"`
		Page       int                       `json:"page"`
		PageSize   int                       `json:"page_size"`
		TotalPages int                       `json:"total_pages"`
	}

	items := make([]dto.LedgerEntryResponse, len(result.Items))
	for i, e := range result.Items {
		items[i] = dto.LedgerEntryFromEntity(&e)
	}

	httputil.JSON(w, http.StatusOK, pagedResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *LedgerHandler) Credit(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	var req dto.CreditRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := h.uc.Credit(r.Context(), franchiseID, req.Amount, req.Description, req.ReferenceType, req.ReferenceID); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "ledger not found")
			return
		}
		httputil.InternalError(w, "failed to credit ledger")
		return
	}
	httputil.NoContent(w)
}
