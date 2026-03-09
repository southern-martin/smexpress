package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
	"github.com/smexpress/services/customer-service/internal/interface/http/dto"
	"github.com/smexpress/services/customer-service/internal/usecase"
)

type CustomerHandler struct {
	uc *usecase.CustomerUseCase
}

func NewCustomerHandler(uc *usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{uc: uc}
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	franchiseID := httputil.QueryString(r, "franchise_id", "")
	search := httputil.QueryString(r, "search", "")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.List(r.Context(), countryCode, franchiseID, search, page)
	if err != nil {
		httputil.InternalError(w, "failed to list customers")
		return
	}

	type pagedResponse struct {
		Items      []dto.CustomerResponse `json:"items"`
		TotalCount int64                  `json:"total_count"`
		Page       int                    `json:"page"`
		PageSize   int                    `json:"page_size"`
		TotalPages int                    `json:"total_pages"`
	}

	items := make([]dto.CustomerResponse, len(result.Items))
	for i, c := range result.Items {
		items[i] = dto.CustomerFromEntity(&c)
	}

	httputil.JSON(w, http.StatusOK, pagedResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *CustomerHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "customer not found")
			return
		}
		httputil.InternalError(w, "failed to get customer")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.CustomerFromEntity(c))
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCustomerRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	paymentTerms := req.PaymentTerms
	if paymentTerms == 0 {
		paymentTerms = 30
	}

	c := &entity.Customer{
		CountryCode:   req.CountryCode,
		FranchiseID:   req.FranchiseID,
		CompanyName:   req.CompanyName,
		TradingName:   req.TradingName,
		AccountNumber: req.AccountNumber,
		ABN:           req.ABN,
		Email:         req.Email,
		Phone:         req.Phone,
		Website:       req.Website,
		CreditLimit:   req.CreditLimit,
		PaymentTerms:  paymentTerms,
	}

	if err := h.uc.Create(r.Context(), c); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "customer already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create customer")
		return
	}
	httputil.Created(w, dto.CustomerFromEntity(c))
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateCustomerRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	existing, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "customer not found")
			return
		}
		httputil.InternalError(w, "failed to get customer")
		return
	}

	if req.CompanyName != "" {
		existing.CompanyName = req.CompanyName
	}
	if req.TradingName != "" {
		existing.TradingName = req.TradingName
	}
	if req.ABN != "" {
		existing.ABN = req.ABN
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.Website != "" {
		existing.Website = req.Website
	}
	if req.CreditLimit != nil {
		existing.CreditLimit = *req.CreditLimit
	}
	if req.PaymentTerms != nil {
		existing.PaymentTerms = *req.PaymentTerms
	}

	if err := h.uc.Update(r.Context(), existing); err != nil {
		httputil.InternalError(w, "failed to update customer")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.CustomerFromEntity(existing))
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "customer not found")
			return
		}
		httputil.InternalError(w, "failed to delete customer")
		return
	}
	httputil.NoContent(w)
}

func (h *CustomerHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Activate(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "customer not found")
			return
		}
		httputil.InternalError(w, "failed to activate customer")
		return
	}
	httputil.NoContent(w)
}

func (h *CustomerHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Deactivate(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "customer not found")
			return
		}
		httputil.InternalError(w, "failed to deactivate customer")
		return
	}
	httputil.NoContent(w)
}

func (h *CustomerHandler) SetCreditHold(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct{ Hold bool `json:"hold"` }
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}
	if err := h.uc.SetCreditHold(r.Context(), id, req.Hold); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "customer not found")
			return
		}
		httputil.InternalError(w, "failed to set credit hold")
		return
	}
	httputil.NoContent(w)
}
