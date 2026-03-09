package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/interface/http/dto"
	"github.com/smexpress/services/config-service/internal/usecase"
)

type HolidayHandler struct {
	uc *usecase.HolidayUseCase
}

func NewHolidayHandler(uc *usecase.HolidayUseCase) *HolidayHandler {
	return &HolidayHandler{uc: uc}
}

func (h *HolidayHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	year := httputil.QueryInt(r, "year", time.Now().Year())
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.List(r.Context(), countryCode, year, page)
	if err != nil {
		httputil.InternalError(w, "failed to list holidays")
		return
	}

	items := make([]dto.HolidayResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = dto.HolidayFromEntity(&item)
	}
	httputil.JSON(w, http.StatusOK, db.PagedResult[dto.HolidayResponse]{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *HolidayHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateHolidayRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	date, err := time.Parse("2006-01-02", req.HolidayDate)
	if err != nil {
		httputil.BadRequest(w, "invalid date format, use YYYY-MM-DD")
		return
	}

	e := &entity.Holiday{
		CountryCode: req.CountryCode,
		HolidayDate: date,
		Name:        req.Name,
		IsRecurring: req.IsRecurring,
	}

	if err := h.uc.Create(r.Context(), e); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "holiday already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create holiday")
		return
	}
	httputil.Created(w, dto.HolidayFromEntity(e))
}

func (h *HolidayHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "holiday not found")
			return
		}
		httputil.InternalError(w, "failed to delete holiday")
		return
	}
	httputil.NoContent(w)
}
