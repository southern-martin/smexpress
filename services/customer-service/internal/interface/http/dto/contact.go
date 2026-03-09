package dto

import (
	"github.com/smexpress/services/customer-service/internal/domain/entity"
)

type CreateContactRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
	Position  string `json:"position"`
	IsPrimary bool   `json:"is_primary"`
	IsBilling bool   `json:"is_billing"`
}

type UpdateContactRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
	Position  string `json:"position"`
	IsPrimary *bool  `json:"is_primary,omitempty"`
	IsBilling *bool  `json:"is_billing,omitempty"`
}

type ContactResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Mobile    string `json:"mobile"`
	Position  string `json:"position"`
	IsPrimary bool   `json:"is_primary"`
	IsBilling bool   `json:"is_billing"`
}

func ContactFromEntity(c *entity.CustomerContact) ContactResponse {
	return ContactResponse{
		ID:        c.ID,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
		Mobile:    c.Mobile,
		Position:  c.Position,
		IsPrimary: c.IsPrimary,
		IsBilling: c.IsBilling,
	}
}
