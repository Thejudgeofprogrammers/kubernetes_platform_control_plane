package domain

import "time"

type APIService struct {
	ID        string
	Name      string
	BaseURL   string
	Protocol  string
	Status    string
	CreatedAt time.Time
}
