package domain

import "time"

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Role      AccessRole `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
}
