package handler

import (
	"net/http"

	"github.com/smexpress/pkg/httputil"
	"github.com/smexpress/services/auth-service/internal/interface/http/dto"
	"github.com/smexpress/services/auth-service/internal/usecase"
)

type PermissionHandler struct {
	permUC *usecase.PermissionUseCase
}

func NewPermissionHandler(permUC *usecase.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{permUC: permUC}
}

func (h *PermissionHandler) List(w http.ResponseWriter, r *http.Request) {
	module := httputil.QueryString(r, "module", "")

	var perms []dto.PermissionResponse
	var err error

	if module != "" {
		entities, e := h.permUC.ListByModule(r.Context(), module)
		err = e
		for _, p := range entities {
			perms = append(perms, dto.PermissionFromEntity(&p))
		}
	} else {
		entities, e := h.permUC.List(r.Context())
		err = e
		for _, p := range entities {
			perms = append(perms, dto.PermissionFromEntity(&p))
		}
	}

	if err != nil {
		httputil.InternalError(w, "failed to list permissions")
		return
	}
	httputil.JSON(w, http.StatusOK, perms)
}
