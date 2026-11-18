package internal

import (
	"errors"
)

var (
	ErrInvalidHashFormat   = errors.New("invalid password hash format.")
	ErrNonExistentUser     = errors.New("user does not exist.")
	ErrInternalServerError = errors.New("error parsing user id.")
)
