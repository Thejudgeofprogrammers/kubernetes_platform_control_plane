package domain

import "time"

type User struct {
	ID        string
	Email     string
	FullName  string
	Role      AccessRole
	CreatedAt time.Time
}
