package usecase

import (
	"context"
	"fmt"

	"github.com/smexpress/pkg/db"
	"github.com/smexpress/pkg/messaging"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
	"github.com/smexpress/services/auth-service/internal/domain/repository"
)

type UserUseCase struct {
	userRepo     repository.UserRepository
	userRoleRepo repository.UserRoleRepository
	publisher    *messaging.Publisher
}

func NewUserUseCase(userRepo repository.UserRepository, userRoleRepo repository.UserRoleRepository, publisher *messaging.Publisher) *UserUseCase {
	return &UserUseCase{userRepo: userRepo, userRoleRepo: userRoleRepo, publisher: publisher}
}

func (uc *UserUseCase) Create(ctx context.Context, user *entity.User, roleIDs []string) error {
	if user.Email == "" || user.FirstName == "" || user.LastName == "" {
		return fmt.Errorf("%w: email, first_name, last_name required", domainerr.ErrInvalidInput)
	}

	hash, err := entity.HashPassword(user.PasswordHash) // PasswordHash field holds plain password at this point
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user.PasswordHash = hash
	user.IsActive = true

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return err
	}

	if len(roleIDs) > 0 {
		if err := uc.userRoleRepo.AssignRoles(ctx, user.ID, roleIDs); err != nil {
			return fmt.Errorf("assign roles: %w", err)
		}
	}

	// Publish user created event for user-service
	if uc.publisher != nil {
		uc.publisher.Publish(ctx, messaging.Event{
			Subject: "auth.user.created",
			Data: map[string]string{
				"user_id":      user.ID,
				"email":        user.Email,
				"first_name":   user.FirstName,
				"last_name":    user.LastName,
				"country_code": user.CountryCode,
				"franchise_id": user.FranchiseID,
			},
		})
	}

	return nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	roles, _ := uc.userRoleRepo.GetUserRoles(ctx, id)
	user.Roles = roles
	return user, nil
}

func (uc *UserUseCase) Update(ctx context.Context, user *entity.User) error {
	return uc.userRepo.Update(ctx, user)
}

func (uc *UserUseCase) Delete(ctx context.Context, id string) error {
	return uc.userRepo.Delete(ctx, id)
}

func (uc *UserUseCase) List(ctx context.Context, countryCode, search string, page db.Page) (db.PagedResult[entity.User], error) {
	return uc.userRepo.List(ctx, countryCode, search, page)
}
