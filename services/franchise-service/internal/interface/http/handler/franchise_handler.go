package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
	"github.com/smexpress/services/franchise-service/internal/interface/http/dto"
	"github.com/smexpress/services/franchise-service/internal/usecase"
)

type FranchiseHandler struct {
	uc *usecase.FranchiseUseCase
}

func NewFranchiseHandler(uc *usecase.FranchiseUseCase) *FranchiseHandler {
	return &FranchiseHandler{uc: uc}
}

func (h *FranchiseHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.List(r.Context(), countryCode, page)
	if err != nil {
		httputil.InternalError(w, "failed to list franchises")
		return
	}

	type pagedResponse struct {
		Items      []dto.FranchiseResponse `json:"items"`
		TotalCount int64                   `json:"total_count"`
		Page       int                     `json:"page"`
		PageSize   int                     `json:"page_size"`
		TotalPages int                     `json:"total_pages"`
	}

	items := make([]dto.FranchiseResponse, len(result.Items))
	for i, f := range result.Items {
		items[i] = dto.FranchiseFromEntity(&f)
	}

	httputil.JSON(w, http.StatusOK, pagedResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *FranchiseHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "franchise not found")
			return
		}
		httputil.InternalError(w, "failed to get franchise")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.FranchiseFromEntity(f))
}

func (h *FranchiseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFranchiseRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	f := &entity.Franchise{
		CountryCode:       req.CountryCode,
		Name:              req.Name,
		Code:              req.Code,
		ContactName:       req.ContactName,
		Email:             req.Email,
		Phone:             req.Phone,
		AddressLine1:      req.AddressLine1,
		AddressLine2:      req.AddressLine2,
		City:              req.City,
		State:             req.State,
		Postcode:          req.Postcode,
		CommissionRate:    req.CommissionRate,
		ParentFranchiseID: req.ParentFranchiseID,
	}

	if err := h.uc.Create(r.Context(), f); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "franchise already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create franchise")
		return
	}
	httputil.Created(w, dto.FranchiseFromEntity(f))
}

func (h *FranchiseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateFranchiseRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "franchise not found")
			return
		}
		httputil.InternalError(w, "failed to get franchise")
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.ContactName != "" {
		existing.ContactName = req.ContactName
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.AddressLine1 != "" {
		existing.AddressLine1 = req.AddressLine1
	}
	if req.AddressLine2 != "" {
		existing.AddressLine2 = req.AddressLine2
	}
	if req.City != "" {
		existing.City = req.City
	}
	if req.State != "" {
		existing.State = req.State
	}
	if req.Postcode != "" {
		existing.Postcode = req.Postcode
	}
	if req.CommissionRate != nil {
		existing.CommissionRate = *req.CommissionRate
	}

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update franchise")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.FranchiseFromEntity(existing))
}

func (h *FranchiseHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Activate(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "franchise not found")
			return
		}
		httputil.InternalError(w, "failed to activate franchise")
		return
	}
	httputil.NoContent(w)
}

func (h *FranchiseHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Deactivate(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "franchise not found")
			return
		}
		httputil.InternalError(w, "failed to deactivate franchise")
		return
	}
	httputil.NoContent(w)
}

func (h *FranchiseHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	settings, err := h.uc.GetSettings(r.Context(), franchiseID)
	if err != nil {
		httputil.InternalError(w, "failed to get settings")
		return
	}
	resp := make([]dto.SettingResponse, len(settings))
	for i, s := range settings {
		resp[i] = dto.SettingFromEntity(&s)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *FranchiseHandler) SetSetting(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	var req dto.SetSettingRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := h.uc.SetSetting(r.Context(), franchiseID, req.Key, req.Value); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to set setting")
		return
	}
	httputil.NoContent(w)
}
