package dto

import (
	"time"

	"github.com/smexpress/services/user-service/internal/domain/entity"
)

type UpdateProfileRequest struct {
	Phone      string `json:"phone"`
	Mobile     string `json:"mobile"`
	JobTitle   string `json:"job_title"`
	Department string `json:"department"`
	AvatarURL  string `json:"avatar_url"`
	Timezone   string `json:"timezone"`
	Locale     string `json:"locale"`
}

type SetPreferenceRequest struct {
	Value string `json:"value"`
}

type ProfileResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CountryCode string    `json:"country_code"`
	Phone       string    `json:"phone,omitempty"`
	Mobile      string    `json:"mobile,omitempty"`
	JobTitle    string    `json:"job_title,omitempty"`
	Department  string    `json:"department,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	Timezone    string    `json:"timezone,omitempty"`
	Locale      string    `json:"locale"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PreferenceResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func ProfileFromEntity(p *entity.UserProfile) ProfileResponse {
	return ProfileResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		CountryCode: p.CountryCode,
		Phone:       p.Phone,
		Mobile:      p.Mobile,
		JobTitle:    p.JobTitle,
		Department:  p.Department,
		AvatarURL:   p.AvatarURL,
		Timezone:    p.Timezone,
		Locale:      p.Locale,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
