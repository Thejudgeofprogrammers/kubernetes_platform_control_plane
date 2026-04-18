package domain

import "time"

type Metric struct {
	ClientID  string
	Requests  int64
	Errors    int64
	Latency   int64
	CreatedAt time.Time
}
