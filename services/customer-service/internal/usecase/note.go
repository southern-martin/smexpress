package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/customer-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/customer-service/internal/domain/errors"
	"github.com/smexpress/services/customer-service/internal/domain/repository"
)

type NoteUseCase struct {
	repo repository.NoteRepository
}

func NewNoteUseCase(repo repository.NoteRepository) *NoteUseCase {
	return &NoteUseCase{repo: repo}
}

func (uc *NoteUseCase) Create(ctx context.Context, note *entity.CustomerNote) error {
	if note.CustomerID == "" || note.Note == "" {
		return fmt.Errorf("%w: customer_id and note required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, note)
}

func (uc *NoteUseCase) ListByCustomer(ctx context.Context, customerID string) ([]entity.CustomerNote, error) {
	return uc.repo.ListByCustomer(ctx, customerID)
}

func (uc *NoteUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
