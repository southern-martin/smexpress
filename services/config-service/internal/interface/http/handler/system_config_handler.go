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

type SystemConfigHandler struct {
	uc *usecase.SystemConfigUseCase
}

func NewSystemConfigHandler(uc *usecase.SystemConfigUseCase) *SystemConfigHandler {
	return &SystemConfigHandler{uc: uc}
}

func (h *SystemConfigHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.List(r.Context(), countryCode, page)
	if err != nil {
		httputil.InternalError(w, "failed to list configs")
		return
	}

	items := make([]dto.SystemConfigResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = dto.SystemConfigFromEntity(&item)
	}
	httputil.JSON(w, http.StatusOK, db.PagedResult[dto.SystemConfigResponse]{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *SystemConfigHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSystemConfigRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	dt := req.DataType
	if dt == "" {
		dt = "string"
	}
	e := &entity.SystemConfig{
		CountryCode: req.CountryCode,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		Description: req.Description,
		DataType:    dt,
	}

	if err := h.uc.Create(r.Context(), e); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "config already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create config")
		return
	}
	httputil.Created(w, dto.SystemConfigFromEntity(e))
}

func (h *SystemConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "config not found")
			return
		}
		httputil.InternalError(w, "failed to get config")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.SystemConfigFromEntity(cfg))
}

func (h *SystemConfigHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateSystemConfigRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "config not found")
			return
		}
		httputil.InternalError(w, "failed to get config")
		return
	}

	existing.ConfigValue = req.ConfigValue
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.DataType != "" {
		existing.DataType = req.DataType
	}

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update config")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.SystemConfigFromEntity(existing))
}

func (h *SystemConfigHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "config not found")
			return
		}
		httputil.InternalError(w, "failed to delete config")
		return
	}
	httputil.NoContent(w)
}
