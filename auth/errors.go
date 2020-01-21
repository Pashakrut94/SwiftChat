package auth

import (
	"github.com/pkg/errors"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrValidationError   = errors.New("validation error")
	ErrUnauthorized      = errors.New("unauthorized  error")
)
