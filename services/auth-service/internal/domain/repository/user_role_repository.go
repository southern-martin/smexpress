package repository

import (
	"context"

	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type UserRoleRepository interface {
	AssignRoles(ctx context.Context, userID string, roleIDs []string) error
	GetUserRoles(ctx context.Context, userID string) ([]entity.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]entity.Permission, error)
}
