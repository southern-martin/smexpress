package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/smexpress/pkg/auth"
	"github.com/smexpress/pkg/messaging"
	"github.com/smexpress/services/auth-service/internal/domain/entity"
	domainerr "github.com/smexpress/services/auth-service/internal/domain/errors"
	"github.com/smexpress/services/auth-service/internal/domain/repository"
)

type LoginResult struct {
	User        *entity.User
	Roles       []entity.Role
	Permissions []entity.Permission
	TokenPair   auth.TokenPair
}

type AuthUseCase struct {
	userRepo         repository.UserRepository
	userRoleRepo     repository.UserRoleRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtCfg           auth.JWTConfig
	publisher        *messaging.Publisher
}

func NewAuthUseCase(
	userRepo repository.UserRepository,
	userRoleRepo repository.UserRoleRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtCfg auth.JWTConfig,
	publisher *messaging.Publisher,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:         userRepo,
		userRoleRepo:     userRoleRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtCfg:           jwtCfg,
		publisher:        publisher,
	}
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, domainerr.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domainerr.ErrAccountInactive
	}
	if user.IsLocked {
		return nil, domainerr.ErrAccountLocked
	}

	if err := user.CheckPassword(password); err != nil {
		user.IncrementFailedAttempts()
		uc.userRepo.Update(ctx, user)
		return nil, domainerr.ErrInvalidCredentials
	}

	user.ResetFailedAttempts()
	now := time.Now()
	user.LastLoginAt = &now
	uc.userRepo.Update(ctx, user)

	roles, _ := uc.userRoleRepo.GetUserRoles(ctx, user.ID)
	permissions, _ := uc.userRoleRepo.GetUserPermissions(ctx, user.ID)

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	claims := auth.Claims{
		UserID:      user.ID,
		Email:       user.Email,
		Roles:       roleNames,
		CountryCode: user.CountryCode,
		FranchiseID: user.FranchiseID,
	}

	tokenPair, err := auth.GenerateTokenPair(uc.jwtCfg, claims)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	// Store refresh token hash
	hash := hashToken(tokenPair.RefreshToken)
	rt := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(uc.jwtCfg.RefreshTokenExpiry),
	}
	uc.refreshTokenRepo.Create(ctx, rt)

	return &LoginResult{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		TokenPair:   tokenPair,
	}, nil
}

func (uc *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	// Parse refresh JWT to get user ID
	claims, err := auth.ParseToken(refreshToken, uc.jwtCfg.SecretKey)
	if err != nil {
		// Try parsing as standard claims (refresh token only has Subject)
		return nil, domainerr.ErrTokenExpired
	}

	hash := hashToken(refreshToken)
	rt, err := uc.refreshTokenRepo.GetByHash(ctx, hash)
	if err != nil {
		return nil, domainerr.ErrInvalidCredentials
	}
	if rt.Revoked {
		return nil, domainerr.ErrTokenRevoked
	}
	if rt.IsExpired() {
		return nil, domainerr.ErrTokenExpired
	}

	// Revoke old token
	uc.refreshTokenRepo.RevokeByUserID(ctx, rt.UserID)

	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, domainerr.ErrNotFound
	}

	roles, _ := uc.userRoleRepo.GetUserRoles(ctx, user.ID)
	permissions, _ := uc.userRoleRepo.GetUserPermissions(ctx, user.ID)

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	newClaims := auth.Claims{
		UserID:      user.ID,
		Email:       user.Email,
		Roles:       roleNames,
		CountryCode: user.CountryCode,
		FranchiseID: user.FranchiseID,
	}

	tokenPair, err := auth.GenerateTokenPair(uc.jwtCfg, newClaims)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	newHash := hashToken(tokenPair.RefreshToken)
	newRT := &entity.RefreshToken{
		UserID:    user.ID,
		TokenHash: newHash,
		ExpiresAt: time.Now().Add(uc.jwtCfg.RefreshTokenExpiry),
	}
	uc.refreshTokenRepo.Create(ctx, newRT)

	return &LoginResult{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		TokenPair:   tokenPair,
	}, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, userID string) error {
	return uc.refreshTokenRepo.RevokeByUserID(ctx, userID)
}

func (uc *AuthUseCase) GetMe(ctx context.Context, userID string) (*entity.User, []entity.Role, []entity.Permission, error) {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	roles, _ := uc.userRoleRepo.GetUserRoles(ctx, user.ID)
	permissions, _ := uc.userRoleRepo.GetUserPermissions(ctx, user.ID)
	return user, roles, permissions, nil
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := user.CheckPassword(oldPassword); err != nil {
		return domainerr.ErrInvalidCredentials
	}

	hash, err := entity.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.PasswordHash = hash
	now := time.Now()
	user.PasswordChangedAt = &now

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Revoke all refresh tokens
	return uc.refreshTokenRepo.RevokeByUserID(ctx, userID)
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
