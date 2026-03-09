package dto

import (
	"github.com/smexpress/services/address-service/internal/domain/entity"
)

type CreateRegionRequest struct {
	CountryCode    string  `json:"country_code"`
	Name           string  `json:"name"`
	Code           string  `json:"code"`
	ParentRegionID *string `json:"parent_region_id,omitempty"`
}

type UpdateRegionRequest struct {
	Name string `json:"name"`
}

type RegionResponse struct {
	ID             string  `json:"id"`
	CountryCode    string  `json:"country_code"`
	Name           string  `json:"name"`
	Code           string  `json:"code"`
	ParentRegionID *string `json:"parent_region_id,omitempty"`
}

func RegionFromEntity(r *entity.Region) RegionResponse {
	return RegionResponse{
		ID:             r.ID,
		CountryCode:    r.CountryCode,
		Name:           r.Name,
		Code:           r.Code,
		ParentRegionID: r.ParentRegionID,
	}
}
