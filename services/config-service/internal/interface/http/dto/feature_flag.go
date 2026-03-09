package dto

import (
	"time"

	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type CreateFeatureFlagRequest struct {
	CountryCode string `json:"country_code"`
	FlagKey     string `json:"flag_key"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type UpdateFeatureFlagRequest struct {
	FlagKey     string `json:"flag_key"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type FeatureFlagResponse struct {
	ID          string    `json:"id"`
	CountryCode string    `json:"country_code"`
	FlagKey     string    `json:"flag_key"`
	Enabled     bool      `json:"enabled"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func FeatureFlagFromEntity(e *entity.FeatureFlag) FeatureFlagResponse {
	return FeatureFlagResponse{
		ID:          e.ID,
		CountryCode: e.CountryCode,
		FlagKey:     e.FlagKey,
		Enabled:     e.Enabled,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
