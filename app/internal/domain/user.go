package domain

import "time"

type User struct {
	ID        string
	Email     string
	FullName  string
	Role      string
	CreatedAt time.Time
}

type AuthClientAccess struct {
	UserID    string
	ClientID  string
	Role      string
	GrantedAt time.Time
}
