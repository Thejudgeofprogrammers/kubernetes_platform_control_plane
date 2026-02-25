package domain

import "time"

type EmailCode struct {
	Email      string
	Code       string
	ExpiresAt time.Time
}
