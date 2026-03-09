package repository

import (
	"context"

	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type PermissionRepository interface {
	List(ctx context.Context) ([]entity.Permission, error)
	ListByModule(ctx context.Context, module string) ([]entity.Permission, error)
}
