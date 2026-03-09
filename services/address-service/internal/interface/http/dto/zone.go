package dto

import (
	"github.com/smexpress/services/address-service/internal/domain/entity"
)

type CreateZoneRequest struct {
	CountryCode string                `json:"country_code"`
	ZoneName    string                `json:"zone_name"`
	ZoneCode    string                `json:"zone_code"`
	Description string                `json:"description"`
	Postcodes   []ZonePostcodeRequest `json:"postcodes,omitempty"`
}

type UpdateZoneRequest struct {
	ZoneName    string `json:"zone_name"`
	Description string `json:"description"`
}

type ZonePostcodeRequest struct {
	PostcodeFrom string `json:"postcode_from"`
	PostcodeTo   string `json:"postcode_to"`
}

type SetPostcodesRequest struct {
	Postcodes []ZonePostcodeRequest `json:"postcodes"`
}

type ZoneResponse struct {
	ID          string                 `json:"id"`
	CountryCode string                 `json:"country_code"`
	ZoneName    string                 `json:"zone_name"`
	ZoneCode    string                 `json:"zone_code"`
	Description string                 `json:"description"`
	Postcodes   []ZonePostcodeResponse `json:"postcodes,omitempty"`
}

type ZonePostcodeResponse struct {
	ID           string `json:"id"`
	PostcodeFrom string `json:"postcode_from"`
	PostcodeTo   string `json:"postcode_to"`
}

func ZoneFromEntity(z *entity.Zone) ZoneResponse {
	postcodes := make([]ZonePostcodeResponse, len(z.Postcodes))
	for i, p := range z.Postcodes {
		postcodes[i] = ZonePostcodeResponse{
			ID:           p.ID,
			PostcodeFrom: p.PostcodeFrom,
			PostcodeTo:   p.PostcodeTo,
		}
	}
	return ZoneResponse{
		ID:          z.ID,
		CountryCode: z.CountryCode,
		ZoneName:    z.ZoneName,
		ZoneCode:    z.ZoneCode,
		Description: z.Description,
		Postcodes:   postcodes,
	}
}
