package handler

import (
	"encoding/json"
	"go_dfs/internal/service"
	"net/http"
)
type HealthHandler struct{
	Service *service.HealthService
}

func NewHealthHandler(svc *service.HealthService) *HealthHandler{
	return &HealthHandler{
		Service:svc,
	}
}

func(h *HealthHandler) Health(w http.ResponseWriter, r *http.Request){
	status :=h.Service.Check()
	response:=map[string]string{
		"status":status,
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(response)
}