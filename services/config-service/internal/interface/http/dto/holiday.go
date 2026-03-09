package dto

import (
	"time"

	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type CreateHolidayRequest struct {
	CountryCode string `json:"country_code"`
	HolidayDate string `json:"holiday_date"`
	Name        string `json:"name"`
	IsRecurring bool   `json:"is_recurring"`
}

type HolidayResponse struct {
	ID          string    `json:"id"`
	CountryCode string    `json:"country_code"`
	HolidayDate string    `json:"holiday_date"`
	Name        string    `json:"name"`
	IsRecurring bool      `json:"is_recurring"`
	CreatedAt   time.Time `json:"created_at"`
}

func HolidayFromEntity(e *entity.Holiday) HolidayResponse {
	return HolidayResponse{
		ID:          e.ID,
		CountryCode: e.CountryCode,
		HolidayDate: e.HolidayDate.Format("2006-01-02"),
		Name:        e.Name,
		IsRecurring: e.IsRecurring,
		CreatedAt:   e.CreatedAt,
	}
}
