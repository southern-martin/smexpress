package dto

import (
	"github.com/smexpress/services/address-service/internal/domain/entity"
)

type CreatePostcodeRequest struct {
	CountryCode string   `json:"country_code"`
	Postcode    string   `json:"postcode"`
	Suburb      string   `json:"suburb"`
	City        string   `json:"city"`
	State       string   `json:"state"`
	StateCode   string   `json:"state_code"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}

type PostcodeResponse struct {
	ID          string   `json:"id"`
	CountryCode string   `json:"country_code"`
	Postcode    string   `json:"postcode"`
	Suburb      string   `json:"suburb"`
	City        string   `json:"city"`
	State       string   `json:"state"`
	StateCode   string   `json:"state_code"`
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
}

func PostcodeFromEntity(p *entity.Postcode) PostcodeResponse {
	return PostcodeResponse{
		ID:          p.ID,
		CountryCode: p.CountryCode,
		Postcode:    p.Postcode,
		Suburb:      p.Suburb,
		City:        p.City,
		State:       p.State,
		StateCode:   p.StateCode,
		Latitude:    p.Latitude,
		Longitude:   p.Longitude,
	}
}
