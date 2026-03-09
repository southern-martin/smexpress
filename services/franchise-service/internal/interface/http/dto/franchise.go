package dto

import (
	"time"

	"github.com/smexpress/services/franchise-service/internal/domain/entity"
)

type CreateFranchiseRequest struct {
	CountryCode       string  `json:"country_code"`
	Name              string  `json:"name"`
	Code              string  `json:"code"`
	ContactName       string  `json:"contact_name"`
	Email             string  `json:"email"`
	Phone             string  `json:"phone"`
	AddressLine1      string  `json:"address_line1"`
	AddressLine2      string  `json:"address_line2"`
	City              string  `json:"city"`
	State             string  `json:"state"`
	Postcode          string  `json:"postcode"`
	CommissionRate    float64 `json:"commission_rate"`
	ParentFranchiseID *string `json:"parent_franchise_id,omitempty"`
}

type UpdateFranchiseRequest struct {
	Name           string   `json:"name"`
	ContactName    string   `json:"contact_name"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	AddressLine1   string   `json:"address_line1"`
	AddressLine2   string   `json:"address_line2"`
	City           string   `json:"city"`
	State          string   `json:"state"`
	Postcode       string   `json:"postcode"`
	CommissionRate *float64 `json:"commission_rate,omitempty"`
}

type SetSettingRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type FranchiseResponse struct {
	ID                string    `json:"id"`
	CountryCode       string    `json:"country_code"`
	Name              string    `json:"name"`
	Code              string    `json:"code"`
	ContactName       string    `json:"contact_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	AddressLine1      string    `json:"address_line1"`
	AddressLine2      string    `json:"address_line2"`
	City              string    `json:"city"`
	State             string    `json:"state"`
	Postcode          string    `json:"postcode"`
	IsActive          bool      `json:"is_active"`
	CommissionRate    float64   `json:"commission_rate"`
	ParentFranchiseID *string   `json:"parent_franchise_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func FranchiseFromEntity(f *entity.Franchise) FranchiseResponse {
	return FranchiseResponse{
		ID:                f.ID,
		CountryCode:       f.CountryCode,
		Name:              f.Name,
		Code:              f.Code,
		ContactName:       f.ContactName,
		Email:             f.Email,
		Phone:             f.Phone,
		AddressLine1:      f.AddressLine1,
		AddressLine2:      f.AddressLine2,
		City:              f.City,
		State:             f.State,
		Postcode:          f.Postcode,
		IsActive:          f.IsActive,
		CommissionRate:    f.CommissionRate,
		ParentFranchiseID: f.ParentFranchiseID,
		CreatedAt:         f.CreatedAt,
		UpdatedAt:         f.UpdatedAt,
	}
}

type SettingResponse struct {
	ID       string `json:"id"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

func SettingFromEntity(s *entity.FranchiseSetting) SettingResponse {
	return SettingResponse{ID: s.ID, Key: s.SettingKey, Value: s.SettingValue}
}
