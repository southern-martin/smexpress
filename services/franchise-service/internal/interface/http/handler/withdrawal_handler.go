package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/auth"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/franchise-service/internal/domain/errors"
	"github.com/smexpress/services/franchise-service/internal/interface/http/dto"
	"github.com/smexpress/services/franchise-service/internal/usecase"
)

type WithdrawalHandler struct {
	uc *usecase.WithdrawalUseCase
}

func NewWithdrawalHandler(uc *usecase.WithdrawalUseCase) *WithdrawalHandler {
	return &WithdrawalHandler{uc: uc}
}

func (h *WithdrawalHandler) Request(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWithdrawalRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	claims, _ := auth.GetClaims(r.Context())
	requestedBy := ""
	if claims != nil {
		requestedBy = claims.UserID
	}

	withdrawal := &entity.FranchiseWithdrawal{
		FranchiseID:       req.FranchiseID,
		CountryCode:       req.CountryCode,
		Amount:            req.Amount,
		Currency:          req.Currency,
		RequestedBy:       requestedBy,
		BankAccountName:   req.BankAccountName,
		BankAccountNumber: req.BankAccountNumber,
		BankBSB:           req.BankBSB,
		Notes:             req.Notes,
	}

	if err := h.uc.Request(r.Context(), withdrawal); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		if errors.Is(err, domainerr.ErrInsufficientBalance) {
			httputil.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create withdrawal")
		return
	}
	httputil.Created(w, dto.WithdrawalFromEntity(withdrawal))
}

func (h *WithdrawalHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("withdrawalId")
	withdrawal, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "withdrawal not found")
			return
		}
		httputil.InternalError(w, "failed to get withdrawal")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.WithdrawalFromEntity(withdrawal))
}

func (h *WithdrawalHandler) Approve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("withdrawalId")
	claims, _ := auth.GetClaims(r.Context())
	approvedBy := ""
	if claims != nil {
		approvedBy = claims.UserID
	}

	if err := h.uc.Approve(r.Context(), id, approvedBy); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "withdrawal not found")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to approve withdrawal")
		return
	}
	httputil.NoContent(w)
}

func (h *WithdrawalHandler) Reject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("withdrawalId")
	claims, _ := auth.GetClaims(r.Context())
	approvedBy := ""
	if claims != nil {
		approvedBy = claims.UserID
	}

	if err := h.uc.Reject(r.Context(), id, approvedBy); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "withdrawal not found")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to reject withdrawal")
		return
	}
	httputil.NoContent(w)
}

func (h *WithdrawalHandler) Process(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("withdrawalId")
	if err := h.uc.Process(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "withdrawal not found")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		if errors.Is(err, domainerr.ErrInsufficientBalance) {
			httputil.Error(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		httputil.InternalError(w, "failed to process withdrawal")
		return
	}
	httputil.NoContent(w)
}

func (h *WithdrawalHandler) ListByFranchise(w http.ResponseWriter, r *http.Request) {
	franchiseID := r.PathValue("id")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.ListByFranchise(r.Context(), franchiseID, page)
	if err != nil {
		httputil.InternalError(w, "failed to list withdrawals")
		return
	}

	type pagedResponse struct {
		Items      []dto.WithdrawalResponse `json:"items"`
		TotalCount int64                    `json:"total_count"`
		Page       int                      `json:"page"`
		PageSize   int                      `json:"page_size"`
		TotalPages int                      `json:"total_pages"`
	}

	items := make([]dto.WithdrawalResponse, len(result.Items))
	for i, w := range result.Items {
		items[i] = dto.WithdrawalFromEntity(&w)
	}

	httputil.JSON(w, http.StatusOK, pagedResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *WithdrawalHandler) ListPending(w http.ResponseWriter, r *http.Request) {
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.ListPending(r.Context(), page)
	if err != nil {
		httputil.InternalError(w, "failed to list pending withdrawals")
		return
	}

	type pagedResponse struct {
		Items      []dto.WithdrawalResponse `json:"items"`
		TotalCount int64                    `json:"total_count"`
		Page       int                      `json:"page"`
		PageSize   int                      `json:"page_size"`
		TotalPages int                      `json:"total_pages"`
	}

	items := make([]dto.WithdrawalResponse, len(result.Items))
	for i, w := range result.Items {
		items[i] = dto.WithdrawalFromEntity(&w)
	}

	httputil.JSON(w, http.StatusOK, pagedResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}
