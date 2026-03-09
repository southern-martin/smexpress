package dto

import (
	"time"

	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type CreateSequenceRequest struct {
	CountryCode   string `json:"country_code"`
	SequenceType  string `json:"sequence_type"`
	Prefix        string `json:"prefix"`
	FormatPattern string `json:"format_pattern"`
}

type SequenceResponse struct {
	ID            string    `json:"id"`
	CountryCode   string    `json:"country_code"`
	SequenceType  string    `json:"sequence_type"`
	Prefix        string    `json:"prefix"`
	CurrentValue  int64     `json:"current_value"`
	FormatPattern string    `json:"format_pattern"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type NextValueResponse struct {
	Value string `json:"value"`
}

func SequenceFromEntity(e *entity.Sequence) SequenceResponse {
	return SequenceResponse{
		ID:            e.ID,
		CountryCode:   e.CountryCode,
		SequenceType:  e.SequenceType,
		Prefix:        e.Prefix,
		CurrentValue:  e.CurrentValue,
		FormatPattern: e.FormatPattern,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
}
