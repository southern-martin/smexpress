package repository

import (
	"context"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, countryCode, search string, page db.Page) (db.PagedResult[entity.User], error)
}
