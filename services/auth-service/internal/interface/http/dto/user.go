package dto

import (
	"time"

	"github.com/smexpress/services/auth-service/internal/domain/entity"
)

type CreateUserRequest struct {
	Email       string   `json:"email"`
	Password    string   `json:"password"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	CountryCode string   `json:"country_code"`
	FranchiseID string   `json:"franchise_id,omitempty"`
	RoleIDs     []string `json:"role_ids"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

type AssignRolesRequest struct {
	RoleIDs []string `json:"role_ids"`
}

type UserResponse struct {
	ID          string         `json:"id"`
	Email       string         `json:"email"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	CountryCode string         `json:"country_code"`
	FranchiseID string         `json:"franchise_id,omitempty"`
	IsActive    bool           `json:"is_active"`
	IsLocked    bool           `json:"is_locked"`
	Roles       []RoleResponse `json:"roles,omitempty"`
	LastLoginAt *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

func UserFromEntity(u *entity.User) UserResponse {
	resp := UserResponse{
		ID:          u.ID,
		Email:       u.Email,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		CountryCode: u.CountryCode,
		FranchiseID: u.FranchiseID,
		IsActive:    u.IsActive,
		IsLocked:    u.IsLocked,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
	}
	for _, r := range u.Roles {
		resp.Roles = append(resp.Roles, RoleFromEntity(&r))
	}
	return resp
}
