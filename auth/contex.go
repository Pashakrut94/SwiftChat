package auth

import (
	"context"
)

const key = "session"

func WithSession(parent context.Context, sess *Session) context.Context {
	// kak rabotaet withValue?
	return context.WithValue(parent, key, sess)
}

func SessionValue(ctx context.Context) *Session {
	v := ctx.Value(key)
	session := v.(*Session)
	return session
}
