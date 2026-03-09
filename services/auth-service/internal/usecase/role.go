package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
	"github.com/smexpress/services/auth-service/internal/domain/repository"
)

type RoleUseCase struct {
	roleRepo repository.RoleRepository
}

func NewRoleUseCase(roleRepo repository.RoleRepository) *RoleUseCase {
	return &RoleUseCase{roleRepo: roleRepo}
}

func (uc *RoleUseCase) Create(ctx context.Context, role *entity.Role) error {
	if role.Name == "" || role.DisplayName == "" {
		return fmt.Errorf("%w: name and display_name required", domainerr.ErrInvalidInput)
	}
	return uc.roleRepo.Create(ctx, role)
}

func (uc *RoleUseCase) GetByID(ctx context.Context, id string) (*entity.Role, error) {
	role, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	perms, _ := uc.roleRepo.GetPermissionsByRoleID(ctx, id)
	role.Permissions = perms
	return role, nil
}

func (uc *RoleUseCase) Update(ctx context.Context, role *entity.Role) error {
	return uc.roleRepo.Update(ctx, role)
}

func (uc *RoleUseCase) Delete(ctx context.Context, id string) error {
	existing, err := uc.roleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.IsSystem {
		return fmt.Errorf("%w: cannot delete system role", domainerr.ErrInvalidInput)
	}
	return uc.roleRepo.Delete(ctx, id)
}

func (uc *RoleUseCase) List(ctx context.Context, countryCode string) ([]entity.Role, error) {
	return uc.roleRepo.List(ctx, countryCode)
}

func (uc *RoleUseCase) SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return uc.roleRepo.SetPermissions(ctx, roleID, permissionIDs)
}
