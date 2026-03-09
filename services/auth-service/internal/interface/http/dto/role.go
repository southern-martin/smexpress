package dto

import "github.com/smexpress/services/auth-service/internal/domain/entity"

type CreateRoleRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	CountryCode string `json:"country_code"`
}

type UpdateRoleRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

type SetPermissionsRequest struct {
	PermissionIDs []string `json:"permission_ids"`
}

type RoleResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	DisplayName string               `json:"display_name"`
	Description string               `json:"description,omitempty"`
	IsSystem    bool                 `json:"is_system"`
	CountryCode string               `json:"country_code"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
}

func RoleFromEntity(r *entity.Role) RoleResponse {
	resp := RoleResponse{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: r.Description,
		IsSystem:    r.IsSystem,
		CountryCode: r.CountryCode,
	}
	for _, p := range r.Permissions {
		resp.Permissions = append(resp.Permissions, PermissionFromEntity(&p))
	}
	return resp
}

type PermissionResponse struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Module string `json:"module"`
}

func PermissionFromEntity(p *entity.Permission) PermissionResponse {
	return PermissionResponse{
		ID:     p.ID,
		Code:   p.Code,
		Name:   p.Name,
		Module: p.Module,
	}
}
