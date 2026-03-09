package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/auth"
	"github.com/smexpress/pkg/httputil"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
	"github.com/smexpress/services/auth-service/internal/interface/http/dto"
	"github.com/smexpress/services/auth-service/internal/usecase"
)

type AuthHandler struct {
	authUC *usecase.AuthUseCase
}

func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	result, err := h.authUC.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domainerr.ErrInvalidCredentials) {
			httputil.Error(w, http.StatusUnauthorized, "invalid email or password")
			return
		}
		if errors.Is(err, domainerr.ErrAccountLocked) {
			httputil.Error(w, http.StatusForbidden, "account is locked")
			return
		}
		if errors.Is(err, domainerr.ErrAccountInactive) {
			httputil.Error(w, http.StatusForbidden, "account is inactive")
			return
		}
		httputil.InternalError(w, "login failed")
		return
	}

	roles := make([]dto.RoleResponse, len(result.Roles))
	for i, r := range result.Roles {
		roles[i] = dto.RoleFromEntity(&r)
	}

	userResp := dto.UserFromEntity(result.User)
	userResp.Roles = roles

	httputil.JSON(w, http.StatusOK, dto.LoginResponse{
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresAt:    result.TokenPair.ExpiresAt,
		User:         userResp,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	result, err := h.authUC.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		httputil.Error(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	httputil.JSON(w, http.StatusOK, dto.LoginResponse{
		AccessToken:  result.TokenPair.AccessToken,
		RefreshToken: result.TokenPair.RefreshToken,
		ExpiresAt:    result.TokenPair.ExpiresAt,
		User:         dto.UserFromEntity(result.User),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	if err := h.authUC.Logout(r.Context(), claims.UserID); err != nil {
		httputil.InternalError(w, "logout failed")
		return
	}
	httputil.NoContent(w)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	user, roles, permissions, err := h.authUC.GetMe(r.Context(), claims.UserID)
	if err != nil {
		httputil.InternalError(w, "failed to get user")
		return
	}

	userResp := dto.UserFromEntity(user)
	for _, r := range roles {
		userResp.Roles = append(userResp.Roles, dto.RoleFromEntity(&r))
	}

	permResp := make([]dto.PermissionResponse, len(permissions))
	for i, p := range permissions {
		permResp[i] = dto.PermissionFromEntity(&p)
	}

	httputil.JSON(w, http.StatusOK, map[string]any{
		"user":        userResp,
		"permissions": permResp,
	})
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.ChangePasswordRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := h.authUC.ChangePassword(r.Context(), claims.UserID, req.OldPassword, req.NewPassword); err != nil {
		if errors.Is(err, domainerr.ErrInvalidCredentials) {
			httputil.Error(w, http.StatusUnauthorized, "incorrect current password")
			return
		}
		httputil.InternalError(w, "failed to change password")
		return
	}
	httputil.JSON(w, http.StatusOK, map[string]string{"message": "password changed"})
}
