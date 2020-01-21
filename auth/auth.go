package auth

import (
	"time"
)

const sessionName = "session-id"

type Session struct {
	SessionID string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
	DeletedAt *time.Time
}
