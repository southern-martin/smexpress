package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/interface/http/dto"
	"github.com/smexpress/services/config-service/internal/usecase"
)

type SequenceHandler struct {
	uc *usecase.SequenceUseCase
}

func NewSequenceHandler(uc *usecase.SequenceUseCase) *SequenceHandler {
	return &SequenceHandler{uc: uc}
}

func (h *SequenceHandler) List(w http.ResponseWriter, r *http.Request) {
	countryCode := httputil.QueryString(r, "country_code", "")
	sequences, err := h.uc.List(r.Context(), countryCode)
	if err != nil {
		httputil.InternalError(w, "failed to list sequences")
		return
	}

	resp := make([]dto.SequenceResponse, len(sequences))
	for i, s := range sequences {
		resp[i] = dto.SequenceFromEntity(&s)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *SequenceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateSequenceRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	pattern := req.FormatPattern
	if pattern == "" {
		pattern = "{prefix}{value}"
	}

	e := &entity.Sequence{
		CountryCode:   req.CountryCode,
		SequenceType:  req.SequenceType,
		Prefix:        req.Prefix,
		FormatPattern: pattern,
	}

	if err := h.uc.Create(r.Context(), e); err != nil {
		if errors.Is(err, domainerr.ErrAlreadyExists) {
			httputil.Error(w, http.StatusConflict, "sequence already exists")
			return
		}
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create sequence")
		return
	}
	httputil.Created(w, dto.SequenceFromEntity(e))
}

func (h *SequenceHandler) NextValue(w http.ResponseWriter, r *http.Request) {
	seqType := r.PathValue("type")
	countryCode := httputil.QueryString(r, "country_code", "")

	value, err := h.uc.NextValue(r.Context(), countryCode, seqType)
	if err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "sequence not found")
			return
		}
		httputil.InternalError(w, "failed to get next value")
		return
	}
	httputil.JSON(w, http.StatusOK, dto.NextValueResponse{Value: value})
}
