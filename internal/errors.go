package internal

import (
	"errors"
)

var (
	ErrInvalidSessionData   = errors.New("invalid session data.")
	ErrInvalidSessionExpiry = errors.New("invalid session expiry date.")
	ErrSessionNotFound      = errors.New("session not found.")
	ErrSessionExpired       = errors.New("session has expired.")
	ErrInvalidHashFormat    = errors.New("invalid password hash format.")
	ErrNonExistentUser      = errors.New("user does not exist.")
	ErrInternalServerError  = errors.New("error parsing user id.")
)
