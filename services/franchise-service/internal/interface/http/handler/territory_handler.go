package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
	"github.com/smexpress/services/franchise-service/internal/interface/http/dto"
	"github.com/smexpress/services/franchise-service/internal/usecase"
)

type TerritoryHandler struct {
	uc *usecase.TerritoryUseCase
}

func NewTerritoryHandler(uc *usecase.TerritoryUseCase) *TerritoryHandler {
	return &TerritoryHandler{uc: uc}
}

func (h *TerritoryHandler) ListByFranchise(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	territories, err := h.uc.ListByFranchise(r.Context(), franchiseID)
	if err != nil {
		httputil.InternalError(w, "failed to list territories")
		return
	}
	resp := make([]dto.TerritoryResponse, len(territories))
	for i, t := range territories {
		resp[i] = dto.TerritoryFromEntity(&t)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *TerritoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("territoryId")
	t, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "territory not found")
			return
		}
		httputil.InternalError(w, "failed to get territory")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.TerritoryFromEntity(t))
}

func (h *TerritoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	var req dto.CreateTerritoryRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	t := &entity.Territory{
		FranchiseID:  franchiseID,
		CountryCode:  req.CountryCode,
		Name:         req.Name,
		PostcodeFrom: req.PostcodeFrom,
		PostcodeTo:   req.PostcodeTo,
		State:        req.State,
		IsExclusive:  req.IsExclusive,
	}

	if err := h.uc.Create(r.Context(), t); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create territory")
		return
	}
	httputil.Created(w, dto.TerritoryFromEntity(t))
}

func (h *TerritoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("territoryId")
	var req dto.UpdateTerritoryRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "territory not found")
			return
		}
		httputil.InternalError(w, "failed to get territory")
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.PostcodeFrom != "" {
		existing.PostcodeFrom = req.PostcodeFrom
	}
	if req.PostcodeTo != "" {
		existing.PostcodeTo = req.PostcodeTo
	}
	if req.State != "" {
		existing.State = req.State
	}
	if req.IsExclusive != nil {
		existing.IsExclusive = *req.IsExclusive
	}

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update territory")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.TerritoryFromEntity(existing))
}

func (h *TerritoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("territoryId")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "territory not found")
			return
		}
		httputil.InternalError(w, "failed to delete territory")
		return
	}
	httputil.NoContent(w)
}

func (h *TerritoryHandler) FindByPostcode(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	postcode := httputil.QueryString(r, "postcode", "")

	territories, err := h.uc.FindByPostcode(r.Context(), countryCode, postcode)
	if err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to find territories")
		return
	}
	resp := make([]dto.TerritoryResponse, len(territories))
	for i, t := range territories {
		resp[i] = dto.TerritoryFromEntity(&t)
	}
	httputil.JSON(w, http.StatusOK, resp)
}
