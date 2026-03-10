package domain

import "time"

type AuthClientAccess struct {
	UserID    string
	ClientID  string
	Role      AccessRole
	GrantedAt time.Time
}
