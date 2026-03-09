package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/interface/http/dto"
	"github.com/smexpress/services/config-service/internal/usecase"
)

type FeatureFlagHandler struct {
	uc *usecase.FeatureFlagUseCase
}

func NewFeatureFlagHandler(uc *usecase.FeatureFlagUseCase) *FeatureFlagHandler {
	return &FeatureFlagHandler{uc: uc}
}

func (h *FeatureFlagHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.List(r.Context(), countryCode, page)
	if err != nil {
		httputil.InternalError(w, "failed to list flags")
		return
	}

	items := make([]dto.FeatureFlagResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = dto.FeatureFlagFromEntity(&item)
	}
	httputil.JSON(w, http.StatusOK, db.PagedResult[dto.FeatureFlagResponse]{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *FeatureFlagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFeatureFlagRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	e := &entity.FeatureFlag{
		CountryCode: req.CountryCode,
		FlagKey:     req.FlagKey,
		Enabled:     req.Enabled,
		Description: req.Description,
	}

	if err := h.uc.Create(r.Context(), e); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "flag already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create flag")
		return
	}
	httputil.Created(w, dto.FeatureFlagFromEntity(e))
}

func (h *FeatureFlagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateFeatureFlagRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "flag not found")
			return
		}
		httputil.InternalError(w, "failed to get flag")
		return
	}

	existing.FlagKey = req.FlagKey
	existing.Enabled = req.Enabled
	existing.Description = req.Description

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update flag")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.FeatureFlagFromEntity(existing))
}

func (h *FeatureFlagHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	flag, err := h.uc.Toggle(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "flag not found")
			return
		}
		httputil.InternalError(w, "failed to toggle flag")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.FeatureFlagFromEntity(flag))
}
