package domain

import "errors"

var ErrClientNotFound = errors.New("client not found")
var ErrClientAlredyExist = errors.New("client already exists")
var ErrInvalidStateTransition = errors.New("invalid state transition")
var ErrConfigVersionExists = errors.New("config version already exists")
var ErrConfigNotFound = errors.New("config not found")
var ErrInvalidRefreshToken = errors.New("invalid refresh token")

var ErrInvalidRole = errors.New("invalid role")
var ErrUserNotFound = errors.New("user not found")
var ErrUnAuthorizedUser = errors.New("unauthorized")
var ErrEmptyUserID = errors.New("empty user id")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrCodeNotFound = errors.New("email code not found")
var ErrCodeExpired = errors.New("email code expired")

var ErrAPIServiceNotFound = errors.New("api service not found")
var ErrDeleteConfig = errors.New("cannot delete active config")
var ErrForbidden = errors.New("forrbiden")