package auth

import (
	"database/sql"
	"time"
)

type SessionRepo struct {
	db *sql.DB
}

func NewSessionRepo(db *sql.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (repo *SessionRepo) Create(sessionID string, userID int, expires time.Time) error {
	q := "insert into sessions (session_id, user_id, created_at, expires_at) values ($1,$2,$3,$4)"
	now := time.Now()
	_, err := repo.db.Exec(q, sessionID, userID, now, expires)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SessionRepo) Delete(sessionID string) error {
	q := "update sessions set deleted_at = $1 where session_id = $2"
	now := time.Now()
	_, err := repo.db.Exec(q, now, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *SessionRepo) get(q string, args ...interface{}) (*Session, error) {
	var s Session
	if err := repo.db.QueryRow(q, args...).Scan(
		&s.SessionID,
		&s.UserID,
		&s.CreatedAt,
		&s.ExpiresAt,
		&s.DeletedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func (repo *SessionRepo) Get(sessionID string) (*Session, error) {
	return repo.get("select session_id, user_id, created_at, expires_at, deleted_at from sessions where session_id = $1 and deleted_at is null", sessionID)
}
