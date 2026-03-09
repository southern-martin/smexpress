package entity

import "time"

type UserProfile struct {
	ID          string
	UserID      string
	CountryCode string
	Phone       string
	Mobile      string
	JobTitle    string
	Department  string
	AvatarURL   string
	Timezone    string
	Locale      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UserPreference struct {
	ID              string
	UserID          string
	PreferenceKey   string
	PreferenceValue string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
