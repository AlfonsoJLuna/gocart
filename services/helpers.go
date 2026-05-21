package services

import (
	"errors"
)

// Error types
var (
	ErrNotFound      = errors.New("not found")
	ErrBadInput      = errors.New("bad input")
	ErrAlreadyExists = errors.New("already exists")
)
