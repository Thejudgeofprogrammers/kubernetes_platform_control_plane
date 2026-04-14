package dto

import (
	"control_plane/internal/domain"
	"time"
)

type APIClientHealth struct {
	ClientID  string              `json:"client_id"`
	Status    domain.HealthStatus `json:"status"`
	LastCheck time.Time           `json:"last_check"`
	Message   string              `json:"message"`
}
