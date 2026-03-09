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

type ZoneHandler struct {
	uc *usecase.ZoneUseCase
}

func NewZoneHandler(uc *usecase.ZoneUseCase) *ZoneHandler {
	return &ZoneHandler{uc: uc}
}

func (h *ZoneHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	zones, err := h.uc.ListByCountry(r.Context(), countryCode)
	if err != nil {
		httputil.InternalError(w, "failed to list zones")
		return
	}
	resp := make([]dto.ZoneResponse, len(zones))
	for i, z := range zones {
		resp[i] = dto.ZoneFromEntity(&z)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *ZoneHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	zone, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "zone not found")
			return
		}
		httputil.InternalError(w, "failed to get zone")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.ZoneFromEntity(zone))
}

func (h *ZoneHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateZoneRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	zone := &entity.Zone{
		CountryCode: req.CountryCode,
		ZoneName:    req.ZoneName,
		ZoneCode:    req.ZoneCode,
		Description: req.Description,
	}

	if err := h.uc.Create(r.Context(), zone); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "zone already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create zone")
		return
	}
	httputil.Created(w, dto.ZoneFromEntity(zone))
}

func (h *ZoneHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateZoneRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "zone not found")
			return
		}
		httputil.InternalError(w, "failed to get zone")
		return
	}

	if req.ZoneName != "" {
		existing.ZoneName = req.ZoneName
	}
	if req.Description != "" {
		existing.Description = req.Description
	}

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update zone")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.ZoneFromEntity(existing))
}

func (h *ZoneHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "zone not found")
			return
		}
		httputil.InternalError(w, "failed to delete zone")
		return
	}
	httputil.NoContent(w)
}

func (h *ZoneHandler) FindZone(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	postcode := httputil.QueryString(r, "postcode", "")

	zone, err := h.uc.FindZone(r.Context(), countryCode, postcode)
	if err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "no zone found for postcode")
			return
		}
		httputil.InternalError(w, "failed to find zone")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.ZoneFromEntity(zone))
}

func (h *ZoneHandler) SetPostcodes(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.SetPostcodesRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	postcodes := make([]entity.ZonePostcode, len(req.Postcodes))
	for i, p := range req.Postcodes {
		postcodes[i] = entity.ZonePostcode{
			ZoneID:       id,
			PostcodeFrom: p.PostcodeFrom,
			PostcodeTo:   p.PostcodeTo,
		}
	}

	if err := h.uc.SetPostcodes(r.Context(), id, postcodes); err != nil {
		httputil.InternalError(w, "failed to set postcodes")
		return
	}
	httputil.NoContent(w)
}
