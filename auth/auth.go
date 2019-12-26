package auth

import (
	"time"
	
	"github.com/nu7hatch/gouuid"
	"net/http"
)



const sessionName = "session-id"

type Session struct {
	SessionID string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
	DeletedAt *time.Time //why pointer?
}

func SetNewSession(sessRepo SessionRepo, userID int, w http.ResponseWriter) error {
	sessID, err := uuid.NewV4()
	if err != nil {		
		return err
	}
	expires := time.Now().Add(time.Minute * 180)
	if err := sessRepo.Create(sessID.String(), userID, expires); err != nil {		
		return err
	}
	// chto takoe httponly? i zachem
	http.SetCookie(w, &http.Cookie{Name:  sessionName, Value: sessID.String(), HttpOnly: true, Expires: expires})
	return nil
}
