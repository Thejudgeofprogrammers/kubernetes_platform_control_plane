package domain

import (
	"time"
)

type AuthType string

const (
	AuthNone   AuthType = "none"
	AuthAPIKey AuthType = "api_key"
	AuthBearer AuthType = "bearer"
)

type APIClientConfig struct {
	ID           string
	ClientID     string
	Version      string
	AuthType     AuthType
	AuthRef      string
	TimeoutMs    int
	RetryCount   int
	RetryBackoff int
	Headers      map[string]string
	CreatedAt    time.Time
	CreatedBy    string
}
