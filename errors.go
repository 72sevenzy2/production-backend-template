package main

import (
	"errors"
)

var (
	ErrNonExistentUser = errors.New("user does not exist.")
)
