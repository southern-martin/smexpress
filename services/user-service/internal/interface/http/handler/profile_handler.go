package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/auth"
	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	domainerr "github.com/smexpress/services/user-service/internal/domain/errors"
	"github.com/smexpress/services/user-service/internal/interface/http/dto"
	"github.com/smexpress/services/user-service/internal/usecase"
)

type ProfileHandler struct {
	uc *usecase.ProfileUseCase
}

func NewProfileHandler(uc *usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{uc: uc}
}

func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	profile, err := h.uc.GetByUserID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "profile not found")
			return
		}
		httputil.InternalError(w, "failed to get profile")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.ProfileFromEntity(profile))
}

func (h *ProfileHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req dto.UpdateProfileRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	profile, err := h.uc.GetByUserID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "profile not found")
			return
		}
		httputil.InternalError(w, "failed to get profile")
		return
	}

	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.Mobile != "" {
		profile.Mobile = req.Mobile
	}
	if req.JobTitle != "" {
		profile.JobTitle = req.JobTitle
	}
	if req.Department != "" {
		profile.Department = req.Department
	}
	if req.AvatarURL != "" {
		profile.AvatarURL = req.AvatarURL
	}
	if req.Timezone != "" {
		profile.Timezone = req.Timezone
	}
	if req.Locale != "" {
		profile.Locale = req.Locale
	}

	if err := h.uc.Update(r.Context(), profile); err != nil {
		httputil.InternalError(w, "failed to update profile")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.ProfileFromEntity(profile))
}

func (h *ProfileHandler) GetByUserID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	profile, err := h.uc.GetByUserID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "profile not found")
			return
		}
		httputil.InternalError(w, "failed to get profile")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.ProfileFromEntity(profile))
}

func (h *ProfileHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.List(r.Context(), countryCode, page)
	if err != nil {
		httputil.InternalError(w, "failed to list profiles")
		return
	}

	items := make([]dto.ProfileResponse, len(result.Items))
	for i, p := range result.Items {
		items[i] = dto.ProfileFromEntity(&p)
	}
	httputil.JSON(w, http.StatusOK, db.PagedResult[dto.ProfileResponse]{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *ProfileHandler) GetPreferences(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	prefs, err := h.uc.GetPreferences(r.Context(), claims.UserID)
	if err != nil {
		httputil.InternalError(w, "failed to get preferences")
		return
	}

	resp := make([]dto.PreferenceResponse, len(prefs))
	for i, p := range prefs {
		resp[i] = dto.PreferenceResponse{Key: p.PreferenceKey, Value: p.PreferenceValue}
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *ProfileHandler) SetPreference(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	key := r.PathValue("key")
	var req dto.SetPreferenceRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	if err := h.uc.SetPreference(r.Context(), claims.UserID, key, req.Value); err != nil {
		httputil.InternalError(w, "failed to set preference")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.PreferenceResponse{Key: key, Value: req.Value})
}

func (h *ProfileHandler) DeletePreference(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.GetClaims(r.Context())
	if !ok {
		httputil.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	key := r.PathValue("key")
	if err := h.uc.DeletePreference(r.Context(), claims.UserID, key); err != nil {
		httputil.InternalError(w, "failed to delete preference")
		return
	}
	httputil.NoContent(w)
}
