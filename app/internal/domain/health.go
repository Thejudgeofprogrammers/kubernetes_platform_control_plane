package domain

import "time"

type HealthStatus string

const (
	HealthHealthy   HealthStatus = "healthy"
	HealthDegraded  HealthStatus = "degraded"
	HealthUnhealthy HealthStatus = "unhealthy"
)

type APIClientHealth struct {
	ClientID  string
	Status    HealthStatus
	LastCheck time.Time
	Message   string
}
