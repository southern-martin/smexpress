package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
	"github.com/smexpress/services/auth-service/internal/interface/http/dto"
	"github.com/smexpress/services/auth-service/internal/usecase"
)

type UserHandler struct {
	userUC     *usecase.UserUseCase
	userRoleUC *usecase.UserUseCase
}

func NewUserHandler(userUC *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	search := httputil.QueryString(r, "search", "")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.userUC.List(r.Context(), countryCode, search, page)
	if err != nil {
		httputil.InternalError(w, "failed to list users")
		return
	}

	items := make([]dto.UserResponse, len(result.Items))
	for i, u := range result.Items {
		items[i] = dto.UserFromEntity(&u)
	}
	httputil.JSON(w, http.StatusOK, db.PagedResult[dto.UserResponse]{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	user := &entity.User{
		Email:        req.Email,
		PasswordHash: req.Password, // Will be hashed in use case
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		CountryCode:  req.CountryCode,
		FranchiseID:  req.FranchiseID,
	}

	if err := h.userUC.Create(r.Context(), user, req.RoleIDs); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "user already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create user")
		return
	}
	httputil.Created(w, dto.UserFromEntity(user))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.userUC.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "user not found")
			return
		}
		httputil.InternalError(w, "failed to get user")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.UserFromEntity(user))
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateUserRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	user, err := h.userUC.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "user not found")
			return
		}
		httputil.InternalError(w, "failed to get user")
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := h.userUC.Update(r.Context(), user); err != nil {
		httputil.InternalError(w, "failed to update user")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.UserFromEntity(user))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.userUC.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "user not found")
			return
		}
		httputil.InternalError(w, "failed to delete user")
		return
	}
	httputil.NoContent(w)
}
