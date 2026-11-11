package main

import (
	"errors"
)

var (
	ErrNonExistentUser     = errors.New("user does not exist.")
	ErrInternalServerError = errors.New("error parsing user id.")
)
