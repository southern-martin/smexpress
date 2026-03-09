package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/address-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/address-service/internal/domain/errors"
	"github.com/smexpress/services/address-service/internal/interface/http/dto"
	"github.com/smexpress/services/address-service/internal/usecase"
)

type RegionHandler struct {
	uc *usecase.RegionUseCase
}

func NewRegionHandler(uc *usecase.RegionUseCase) *RegionHandler {
	return &RegionHandler{uc: uc}
}

func (h *RegionHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	regions, err := h.uc.ListByCountry(r.Context(), countryCode)
	if err != nil {
		httputil.InternalError(w, "failed to list regions")
		return
	}
	resp := make([]dto.RegionResponse, len(regions))
	for i, reg := range regions {
		resp[i] = dto.RegionFromEntity(&reg)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *RegionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	region, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "region not found")
			return
		}
		httputil.InternalError(w, "failed to get region")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.RegionFromEntity(region))
}

func (h *RegionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRegionRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	region := &entity.Region{
		CountryCode:    req.CountryCode,
		Name:           req.Name,
		Code:           req.Code,
		ParentRegionID: req.ParentRegionID,
	}

	if err := h.uc.Create(r.Context(), region); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "region already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create region")
		return
	}
	httputil.Created(w, dto.RegionFromEntity(region))
}

func (h *RegionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateRegionRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "region not found")
			return
		}
		httputil.InternalError(w, "failed to get region")
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update region")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.RegionFromEntity(existing))
}

func (h *RegionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "region not found")
			return
		}
		httputil.InternalError(w, "failed to delete region")
		return
	}
	httputil.NoContent(w)
}
