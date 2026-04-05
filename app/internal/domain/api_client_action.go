package domain

import "time"

type ActionStatus string

const (
	ActionPending   ActionStatus = "pending"
	ActionRunning   ActionStatus = "running"
	ActionSuccess   ActionStatus = "success"
	ActionFailed    ActionStatus = "failed"
)

type APIClientAction struct {
	ID        string
	ClientID  string
	UserID    string
	Type      ActionType
	Status    ActionStatus
	Error     *string
	CreatedAt time.Time
}