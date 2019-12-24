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
	DeletedAt *time.Time
}

func SetNewSession(sessRepo SessionRepo, userID int, w http.ResponseWriter) error {
	unicSessID, err := uuid.NewV4()
	if err != nil {		
		return err
	}
	expires := time.Now().Add(time.Minute * 5)
	if err := sessRepo.Create(unicSessID.String(), userID, expires); err != nil {		
		return err
	}
	http.SetCookie(w, &http.Cookie{Name:  sessionName, Value: unicSessID.String(), HttpOnly: true, Expires: expires})
	return nil
}

func getCookieByName(cookie []*http.Cookie, name string) string {
	cookieLen := len(cookie)
	result := ""
	for i := 0; i < cookieLen; i++ {
		if cookie[i].Name == name {
			result = cookie[i].Value
		}
	}
	return result
}
