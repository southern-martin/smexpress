package repository

import (
	"context"

	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, countryCode string) ([]entity.Role, error)
	SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]entity.Permission, error)
}
