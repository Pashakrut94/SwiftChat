package auth

import (
	"net/http"
	"time"

	"github.com/Pashakrut94/SwiftChat/handlers"
	uuid "github.com/nu7hatch/gouuid"
)

const sessionName = "session-id"

type Session struct {
	SessionID string
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
	DeletedAt *time.Time
}

func SetNewSession(w http.ResponseWriter, sessRepo SessionRepo, userID int) {
	sessID, err := uuid.NewV4()
	if err != nil {
		handlers.HandleResponseError(w, "Error getting hash of the session", http.StatusInternalServerError)
		return
	}
	expires := time.Now().Add(time.Minute * 180)
	if err := sessRepo.Create(sessID.String(), userID, expires); err != nil {
		handlers.HandleResponseError(w, "Error creating new session", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: sessionName, Value: sessID.String(), HttpOnly: true, Expires: expires})
	return
}
