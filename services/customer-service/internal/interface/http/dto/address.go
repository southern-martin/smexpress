package dto

import (
	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type CreateAddressRequest struct {
	AddressType  string `json:"address_type"`
	CompanyName  string `json:"company_name"`
	ContactName  string `json:"contact_name"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
	CountryCode  string `json:"country_code"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	IsDefault    bool   `json:"is_default"`
	Instructions string `json:"instructions"`
}

type UpdateAddressRequest struct {
	CompanyName  string `json:"company_name"`
	ContactName  string `json:"contact_name"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
	CountryCode  string `json:"country_code"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	IsDefault    *bool  `json:"is_default,omitempty"`
	Instructions string `json:"instructions"`
}

type AddressResponse struct {
	ID           string `json:"id"`
	AddressType  string `json:"address_type"`
	CompanyName  string `json:"company_name"`
	ContactName  string `json:"contact_name"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
	CountryCode  string `json:"country_code"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	IsDefault    bool   `json:"is_default"`
	Instructions string `json:"instructions"`
}

func AddressFromEntity(a *entity.CustomerAddress) AddressResponse {
	return AddressResponse{
		ID:           a.ID,
		AddressType:  a.AddressType,
		CompanyName:  a.CompanyName,
		ContactName:  a.ContactName,
		AddressLine1: a.AddressLine1,
		AddressLine2: a.AddressLine2,
		City:         a.City,
		State:        a.State,
		Postcode:     a.Postcode,
		CountryCode:  a.CountryCode,
		Phone:        a.Phone,
		Email:        a.Email,
		IsDefault:    a.IsDefault,
		Instructions: a.Instructions,
	}
}
