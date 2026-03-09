package handler

import (
	"errors"
	"net/http"

	"github.com/smexpress/pkg/auth"
	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
	"github.com/smexpress/services/customer-service/internal/interface/http/dto"
	"github.com/smexpress/services/customer-service/internal/usecase"
)

type NoteHandler struct {
	uc *usecase.NoteUseCase
}

func NewNoteHandler(uc *usecase.NoteUseCase) *NoteHandler {
	return &NoteHandler{uc: uc}
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	notes, err := h.uc.ListByCustomer(r.Context(), customerID)
	if err != nil {
		httputil.InternalError(w, "failed to list notes")
		return
	}
	resp := make([]dto.NoteResponse, len(notes))
	for i, n := range notes {
		resp[i] = dto.NoteFromEntity(&n)
	}
	httputil.JSON(w, http.StatusOK, resp)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	var req dto.CreateNoteRequest
	if err := httputil.Decode(r, &req); err != nil {
		httputil.BadRequest(w, err.Error())
		return
	}

	claims, _ := auth.GetClaims(r.Context())
	createdBy := ""
	if claims != nil {
		createdBy = claims.UserID
	}

	note := &entity.CustomerNote{
		CustomerID: customerID,
		Note:       req.Note,
		CreatedBy:  createdBy,
	}

	if err := h.uc.Create(r.Context(), note); err != nil {
		if errors.Is(err, domainerr.ErrInvalidInput) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create note")
		return
	}
	httputil.Created(w, dto.NoteFromEntity(note))
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	noteID := r.PathValue("noteId")
	if err := h.uc.Delete(r.Context(), noteID); err != nil {
		if errors.Is(err, domainerr.ErrNotFound) {
			httputil.NotFound(w, "note not found")
			return
		}
		httputil.InternalError(w, "failed to delete note")
		return
	}
	httputil.NoContent(w)
}
