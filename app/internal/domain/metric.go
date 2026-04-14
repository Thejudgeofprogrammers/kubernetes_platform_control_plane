package domain

import "time"

type Metric struct {
	ClientID  string
	Name      string
	Value     float64
	CreatedAt time.Time
}
