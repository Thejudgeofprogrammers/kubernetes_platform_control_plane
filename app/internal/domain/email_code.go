package domain

import "time"

type EmailCode struct {
	// ID         string
	Email      string
	Code       string
	ExpiresAt time.Time
	// CreatedAt time.Time
}
