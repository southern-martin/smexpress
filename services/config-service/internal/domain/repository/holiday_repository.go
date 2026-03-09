package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type HolidayRepository interface {
	Create(ctx context.Context, holiday *entity.Holiday) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, countryCode string, year int, page db.Page) (db.PagedResult[entity.Holiday], error)
}
