package domain

import "errors"

var ErrClientNotFound = errors.New("client not found")
var ErrInvalidStateTransition = errors.New("invalid state transition")
var ErrConfigVersionExists = errors.New("config version already exists")
var ErrConfigNotFound = errors.New("config not found")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrCodeNotFound = errors.New("email code not found")
var ErrCodeExpired = errors.New("email code expired")