package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type LedgerRepository interface {
	GetByFranchiseID(ctx context.Context, franchiseID string) (*entity.FranchiseLedger, error)
	Create(ctx context.Context, ledger *entity.FranchiseLedger) error
	AddEntry(ctx context.Context, ledgerID string, entry *entity.FranchiseLedgerEntry) error
	ListEntries(ctx context.Context, ledgerID string, page db.Page) (db.PagedResult[entity.FranchiseLedgerEntry], error)
}
