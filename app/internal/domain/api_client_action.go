package domain

import "time"

type APIClientAction struct {
	ID        string
	ClientID  string
	UserID    string
	Type      ActionType
	CreatedAt time.Time
}
