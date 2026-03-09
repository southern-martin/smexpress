package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
	"github.com/smexpress/services/auth-service/internal/interface/http/dto"
	"github.com/smexpress/services/auth-service/internal/usecase"
)

type RoleHandler struct {
	roleUC *usecase.RoleUseCase
}

func NewRoleHandler(roleUC *usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{roleUC: roleUC}
}

func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	roles, err := h.roleUC.List(r.Context(), countryCode)
	if err != nil {
		httputil.InternalError(w, "failed to list roles")
		return
	}
	resp := make([]dto.RoleResponse, len(roles))
	for i, r := range roles {
		resp[i] = dto.RoleFromEntity(&r)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRoleRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	role := &entity.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		CountryCode: req.CountryCode,
	}

	if err := h.roleUC.Create(r.Context(), role); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "role already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create role")
		return
	}
	httputil.Created(w, dto.RoleFromEntity(role))
}

func (h *RoleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	role, err := h.roleUC.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "role not found")
			return
		}
		httputil.InternalError(w, "failed to get role")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.RoleFromEntity(role))
}

func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateRoleRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	role, err := h.roleUC.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "role not found")
			return
		}
		httputil.InternalError(w, "failed to get role")
		return
	}

	if req.DisplayName != "" {
		role.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		role.Description = req.Description
	}

	if err := h.roleUC.Update(r.Context(), role); err != nil {
		httputil.InternalError(w, "failed to update role")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.RoleFromEntity(role))
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.roleUC.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "role not found")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to delete role")
		return
	}
	httputil.NoContent(w)
}

func (h *RoleHandler) SetPermissions(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.SetPermissionsRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := h.roleUC.SetPermissions(r.Context(), id, req.PermissionIDs); err != nil {
		httputil.InternalError(w, "failed to set permissions")
		return
	}

	role, _ := h.roleUC.GetByID(r.Context(), id)
	httputil.JSON(w, http.StatusOK, dto.RoleFromEntity(role))
}
