package chat

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound          = errors.New("can't find credentials")
	ErrRoomAlreadyExists = errors.New("room already exists")
)
