package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
	"github.com/smexpress/services/customer-service/internal/interface/http/dto"
	"github.com/smexpress/services/customer-service/internal/usecase"
)

type AddressHandler struct {
	uc *usecase.AddressUseCase
}

func NewAddressHandler(uc *usecase.AddressUseCase) *AddressHandler {
	return &AddressHandler{uc: uc}
}

func (h *AddressHandler) List(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	addresses, err := h.uc.ListByCustomer(r.Context(), customerID)
	if err != nil {
		httputil.InternalError(w, "failed to list addresses")
		return
	}
	resp := make([]dto.AddressResponse, len(addresses))
	for i, a := range addresses {
		resp[i] = dto.AddressFromEntity(&a)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	var req dto.CreateAddressRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	addrType := req.AddressType
	if addrType == "" {
		addrType = "shipping"
	}

	addr := &entity.CustomerAddress{
		CustomerID:   customerID,
		AddressType:  addrType,
		CompanyName:  req.CompanyName,
		ContactName:  req.ContactName,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		City:         req.City,
		State:        req.State,
		Postcode:     req.Postcode,
		CountryCode:  req.CountryCode,
		Phone:        req.Phone,
		Email:        req.Email,
		IsDefault:    req.IsDefault,
		Instructions: req.Instructions,
	}

	if err := h.uc.Create(r.Context(), addr); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create address")
		return
	}
	httputil.Created(w, dto.AddressFromEntity(addr))
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	addrID := r.PathValue("addressId")
	var req dto.UpdateAddressRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), addrID)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "address not found")
			return
		}
		httputil.InternalError(w, "failed to get address")
		return
	}

	if req.CompanyName != "" { existing.CompanyName = req.CompanyName }
	if req.ContactName != "" { existing.ContactName = req.ContactName }
	if req.AddressLine1 != "" { existing.AddressLine1 = req.AddressLine1 }
	if req.AddressLine2 != "" { existing.AddressLine2 = req.AddressLine2 }
	if req.City != "" { existing.City = req.City }
	if req.State != "" { existing.State = req.State }
	if req.Postcode != "" { existing.Postcode = req.Postcode }
	if req.CountryCode != "" { existing.CountryCode = req.CountryCode }
	if req.Phone != "" { existing.Phone = req.Phone }
	if req.Email != "" { existing.Email = req.Email }
	if req.IsDefault != nil { existing.IsDefault = *req.IsDefault }
	if req.Instructions != "" { existing.Instructions = req.Instructions }

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update address")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.AddressFromEntity(existing))
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	addrID := r.PathValue("addressId")
	if err := h.uc.Delete(r.Context(), addrID); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "address not found")
			return
		}
		httputil.InternalError(w, "failed to delete address")
		return
	}
	httputil.NoContent(w)
}
