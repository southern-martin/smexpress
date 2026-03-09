package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/interface/http/dto"
	"github.com/smexpress/services/config-service/internal/usecase"
)

type CountryConfigHandler struct {
	uc *usecase.CountryConfigUseCase
}

func NewCountryConfigHandler(uc *usecase.CountryConfigUseCase) *CountryConfigHandler {
	return &CountryConfigHandler{uc: uc}
}

func (h *CountryConfigHandler) List(w http.ResponseWriter, r *http.Request) {
	configs, err := h.uc.List(r.Context())
	if err != nil {
		httputil.InternalError(w, "failed to list countries")
		return
	}
	resp := make([]dto.CountryConfigResponse, len(configs))
	for i, c := range configs {
		resp[i] = dto.CountryConfigFromEntity(&c)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *CountryConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	cfg, err := h.uc.GetByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "country not found")
			return
		}
		httputil.InternalError(w, "failed to get country")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.CountryConfigFromEntity(cfg))
}

func (h *CountryConfigHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCountryConfigRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	e := &entity.CountryConfig{
		CountryCode:    req.CountryCode,
		CountryName:    req.CountryName,
		CurrencyCode:   req.CurrencyCode,
		CurrencySymbol: req.CurrencySymbol,
		Timezone:       req.Timezone,
		DateFormat:     req.DateFormat,
		WeightUnit:     req.WeightUnit,
		DimensionUnit:  req.DimensionUnit,
		Locale:         req.Locale,
		IsActive:       true,
	}

	if err := h.uc.Create(r.Context(), e); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "country already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create country")
		return
	}
	httputil.Created(w, dto.CountryConfigFromEntity(e))
}

func (h *CountryConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	var req dto.UpdateCountryConfigRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByCode(r.Context(), code)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "country not found")
			return
		}
		httputil.InternalError(w, "failed to get country")
		return
	}

	if req.CountryName != "" {
		existing.CountryName = req.CountryName
	}
	if req.CurrencyCode != "" {
		existing.CurrencyCode = req.CurrencyCode
	}
	if req.CurrencySymbol != "" {
		existing.CurrencySymbol = req.CurrencySymbol
	}
	if req.Timezone != "" {
		existing.Timezone = req.Timezone
	}
	if req.DateFormat != "" {
		existing.DateFormat = req.DateFormat
	}
	if req.WeightUnit != "" {
		existing.WeightUnit = req.WeightUnit
	}
	if req.DimensionUnit != "" {
		existing.DimensionUnit = req.DimensionUnit
	}
	if req.Locale != "" {
		existing.Locale = req.Locale
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update country")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.CountryConfigFromEntity(existing))
}
