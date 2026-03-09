package dto

import (
	"time"

	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type CreateTerritoryRequest struct {
	FranchiseID  string `json:"franchise_id"`
	CountryCode  string `json:"country_code"`
	Name         string `json:"name"`
	PostcodeFrom string `json:"postcode_from"`
	PostcodeTo   string `json:"postcode_to"`
	State        string `json:"state"`
	IsExclusive  bool   `json:"is_exclusive"`
}

type UpdateTerritoryRequest struct {
	Name         string `json:"name"`
	PostcodeFrom string `json:"postcode_from"`
	PostcodeTo   string `json:"postcode_to"`
	State        string `json:"state"`
	IsExclusive  *bool  `json:"is_exclusive,omitempty"`
}

type TerritoryResponse struct {
	ID           string    `json:"id"`
	FranchiseID  string    `json:"franchise_id"`
	CountryCode  string    `json:"country_code"`
	Name         string    `json:"name"`
	PostcodeFrom string    `json:"postcode_from"`
	PostcodeTo   string    `json:"postcode_to"`
	State        string    `json:"state"`
	IsExclusive  bool      `json:"is_exclusive"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func TerritoryFromEntity(t *entity.Territory) TerritoryResponse {
	return TerritoryResponse{
		ID:           t.ID,
		FranchiseID:  t.FranchiseID,
		CountryCode:  t.CountryCode,
		Name:         t.Name,
		PostcodeFrom: t.PostcodeFrom,
		PostcodeTo:   t.PostcodeTo,
		State:        t.State,
		IsExclusive:  t.IsExclusive,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}
