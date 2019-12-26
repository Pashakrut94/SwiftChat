package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Pashakrut94/SwiftChat/handlers"
)

//zamikanie
func RequireAuthentication(sessRepo SessionRepo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie(sessionName)
			if err != nil {
				handlers.HandleResponseError(w, "Error getting session from request", http.StatusUnauthorized)
				return
			}
			session, err := sessRepo.Get(c.Value)
			if err == sql.ErrNoRows {
				handlers.HandleResponseError(w, "Error session not found", http.StatusUnauthorized)
				return
			}
			if err != nil {
				handlers.HandleResponseError(w, "Error getting session from DB", http.StatusInternalServerError)
				return
			}
			if time.Now().After(session.ExpiresAt) {
				if err := sessRepo.Delete(c.Value); err != nil {
					handlers.HandleResponseError(w, "Error updating deleted_at", http.StatusUnauthorized)
					return
				}
				handlers.HandleResponseError(w, "Session expired", http.StatusUnauthorized)
				return
			}
			ctx := WithSession(r.Context(), session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
