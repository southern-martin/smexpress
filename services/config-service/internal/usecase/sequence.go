package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/config-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/config-service/internal/domain/errors"
	"github.com/smexpress/services/config-service/internal/domain/repository"
)

type SequenceUseCase struct {
	repo repository.SequenceRepository
}

func NewSequenceUseCase(repo repository.SequenceRepository) *SequenceUseCase {
	return &SequenceUseCase{repo: repo}
}

func (uc *SequenceUseCase) Create(ctx context.Context, seq *entity.Sequence) error {
	if seq.CountryCode == "" || seq.SequenceType == "" {
		return fmt.Errorf("%w: country_code and sequence_type required", domainerr.ErrInvalidInput)
	}
	return uc.repo.Create(ctx, seq)
}

func (uc *SequenceUseCase) List(ctx context.Context, countryCode string) ([]entity.Sequence, error) {
	return uc.repo.List(ctx, countryCode)
}

func (uc *SequenceUseCase) NextValue(ctx context.Context, countryCode, seqType string) (string, error) {
	return uc.repo.NextValue(ctx, countryCode, seqType)
}
