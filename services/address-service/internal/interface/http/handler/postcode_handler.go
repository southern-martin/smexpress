package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/address-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/address-service/internal/domain/errors"
	"github.com/smexpress/services/address-service/internal/interface/http/dto"
	"github.com/smexpress/services/address-service/internal/usecase"
)

type PostcodeHandler struct {
	uc *usecase.PostcodeUseCase
}

func NewPostcodeHandler(uc *usecase.PostcodeUseCase) *PostcodeHandler {
	return &PostcodeHandler{uc: uc}
}

func (h *PostcodeHandler) Search(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	query := httputil.QueryString(r, "q", "")
	page := db.Page{
		Number: httputil.QueryInt(r, "page", 1),
		Size:   httputil.QueryInt(r, "page_size", 20),
	}

	result, err := h.uc.Search(r.Context(), countryCode, query, page)
	if err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to search postcodes")
		return
	}

	type pagedResponse struct {
		Items      []dto.PostcodeResponse `json:"items"`
		TotalCount int64                  `json:"total_count"`
		Page       int                    `json:"page"`
		PageSize   int                    `json:"page_size"`
		TotalPages int                    `json:"total_pages"`
	}

	items := make([]dto.PostcodeResponse, len(result.Items))
	for i, p := range result.Items {
		items[i] = dto.PostcodeFromEntity(&p)
	}

	httputil.JSON(w, http.StatusOK, pagedResponse{
		Items:      items,
		TotalCount: result.TotalCount,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	})
}

func (h *PostcodeHandler) Lookup(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	postcode := httputil.QueryString(r, "postcode", "")

	results, err := h.uc.Lookup(r.Context(), countryCode, postcode)
	if err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to lookup postcode")
		return
	}

	resp := make([]dto.PostcodeResponse, len(results))
	for i, p := range results {
		resp[i] = dto.PostcodeFromEntity(&p)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *PostcodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreatePostcodeRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	p := &entity.Postcode{
		CountryCode: req.CountryCode,
		Postcode:    req.Postcode,
		Suburb:      req.Suburb,
		City:        req.City,
		State:       req.State,
		StateCode:   req.StateCode,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
	}

	if err := h.uc.Create(r.Context(), p); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create postcode")
		return
	}
	httputil.Created(w, dto.PostcodeFromEntity(p))
}
