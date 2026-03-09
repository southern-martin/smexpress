package repository

import (
	"context"

	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type SequenceRepository interface {
	Create(ctx context.Context, seq *entity.Sequence) error
	List(ctx context.Context, countryCode string) ([]entity.Sequence, error)
	NextValue(ctx context.Context, countryCode, seqType string) (string, error)
}
