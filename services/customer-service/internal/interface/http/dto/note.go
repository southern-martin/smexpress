package dto

import (
	"time"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type CreateNoteRequest struct {
	Note string `json:"note"`
}

type NoteResponse struct {
	ID        string    `json:"id"`
	Note      string    `json:"note"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func NoteFromEntity(n *entity.CustomerNote) NoteResponse {
	return NoteResponse{
		ID:        n.ID,
		Note:      n.Note,
		CreatedBy: n.CreatedBy,
		CreatedAt: n.CreatedAt,
	}
}
