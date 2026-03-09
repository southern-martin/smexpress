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

type ContactHandler struct {
	uc *usecase.ContactUseCase
}

func NewContactHandler(uc *usecase.ContactUseCase) *ContactHandler {
	return &ContactHandler{uc: uc}
}

func (h *ContactHandler) List(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	contacts, err := h.uc.ListByCustomer(r.Context(), customerID)
	if err != nil {
		httputil.InternalError(w, "failed to list contacts")
		return
	}
	resp := make([]dto.ContactResponse, len(contacts))
	for i, c := range contacts {
		resp[i] = dto.ContactFromEntity(&c)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *ContactHandler) Create(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	var req dto.CreateContactRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	contact := &entity.CustomerContact{
		CustomerID: customerID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		Phone:      req.Phone,
		Mobile:     req.Mobile,
		Position:   req.Position,
		IsPrimary:  req.IsPrimary,
		IsBilling:  req.IsBilling,
	}

	if err := h.uc.Create(r.Context(), contact); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create contact")
		return
	}
	httputil.Created(w, dto.ContactFromEntity(contact))
}

func (h *ContactHandler) Update(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("contactId")
	var req dto.UpdateContactRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), contactID)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "contact not found")
			return
		}
		httputil.InternalError(w, "failed to get contact")
		return
	}

	if req.FirstName != "" { existing.FirstName = req.FirstName }
	if req.LastName != "" { existing.LastName = req.LastName }
	if req.Email != "" { existing.Email = req.Email }
	if req.Phone != "" { existing.Phone = req.Phone }
	if req.Mobile != "" { existing.Mobile = req.Mobile }
	if req.Position != "" { existing.Position = req.Position }
	if req.IsPrimary != nil { existing.IsPrimary = *req.IsPrimary }
	if req.IsBilling != nil { existing.IsBilling = *req.IsBilling }

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update contact")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.ContactFromEntity(existing))
}

func (h *ContactHandler) Delete(w http.ResponseWriter, r *http.Request) {
	contactID := r.PathValue("contactId")
	if err := h.uc.Delete(r.Context(), contactID); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "contact not found")
			return
		}
		httputil.InternalError(w, "failed to delete contact")
		return
	}
	httputil.NoContent(w)
}
