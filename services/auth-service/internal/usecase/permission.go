package usecase

import (
	"context"

	"github.com/smexpress/services/auth-service/internal/domain/entity"
	"github.com/smexpress/services/auth-service/internal/domain/repository"
)

type PermissionUseCase struct {
	repo repository.PermissionRepository
}

func NewPermissionUseCase(repo repository.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{repo: repo}
}

func (uc *PermissionUseCase) List(ctx context.Context) ([]entity.Permission, error) {
	return uc.repo.List(ctx)
}

func (uc *PermissionUseCase) ListByModule(ctx context.Context, module string) ([]entity.Permission, error) {
	return uc.repo.ListByModule(ctx, module)
}
