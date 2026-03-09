package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type WithdrawalRepository interface {
	Create(ctx context.Context, withdrawal *entity.FranchiseWithdrawal) error
	GetByID(ctx context.Context, id string) (*entity.FranchiseWithdrawal, error)
	Update(ctx context.Context, withdrawal *entity.FranchiseWithdrawal) error
	ListByFranchise(ctx context.Context, franchiseID string, page db.Page) (db.PagedResult[entity.FranchiseWithdrawal], error)
	ListPending(ctx context.Context, page db.Page) (db.PagedResult[entity.FranchiseWithdrawal], error)
}
