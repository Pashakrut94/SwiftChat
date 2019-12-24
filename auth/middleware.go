package auth

import (
	"database/sql"
	"net/http"
	"time"
)

func RequireAuthentication(sessRepo SessionRepo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(sessionName)
			if err != nil {
				http.Error(w, "Error getting session from request",
					http.StatusUnauthorized)
				return
			}
			session, err := sessRepo.Get(c.Value)
			if err == sql.ErrNoRows {
				http.Error(w, "Session not found",
					http.StatusUnauthorized)
				return
			}
			if err != nil {
				http.Error(w, "Error getting session from DB",
					http.StatusInternalServerError)
				return
			}
			if time.Now().After(session.ExpiresAt) {
				if err := sessRepo.Delete(c.Value); err != nil {
					http.Error(w, "Error updating deleted_at",
						http.StatusUnauthorized)
					return
				}
				http.Error(w, "Session expired", http.StatusUnauthorized)
				return
			}
			ctx := WithSession(r.Context(), session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
