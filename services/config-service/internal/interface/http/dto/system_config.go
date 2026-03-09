package dto

import (
	"time"

	"github.com/smexpress/services/config-service/internal/domain/entity"
)

type CreateSystemConfigRequest struct {
	CountryCode string `json:"country_code"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Description string `json:"description"`
	DataType    string `json:"data_type"`
}

type UpdateSystemConfigRequest struct {
	ConfigValue string `json:"config_value"`
	Description string `json:"description"`
	DataType    string `json:"data_type"`
}

type SystemConfigResponse struct {
	ID          string    `json:"id"`
	CountryCode string    `json:"country_code"`
	ConfigKey   string    `json:"config_key"`
	ConfigValue string    `json:"config_value"`
	Description string    `json:"description,omitempty"`
	DataType    string    `json:"data_type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func SystemConfigFromEntity(e *entity.SystemConfig) SystemConfigResponse {
	return SystemConfigResponse{
		ID:          e.ID,
		CountryCode: e.CountryCode,
		ConfigKey:   e.ConfigKey,
		ConfigValue: e.ConfigValue,
		Description: e.Description,
		DataType:    e.DataType,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
