package auth

import (
	"net/http"

	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/pkg/errors"
)

func RequireAuthentication(sessRepo SessionRepo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionName)
			if err != nil {
				handlers.HandleResponseError(w, errors.Wrap(err, "error getting cookie from request").Error(), http.StatusUnauthorized)
				return
			}
			session, err := HandleAuthentication(sessRepo, cookie.Value)
			if err != nil {
				switch errors.Cause(err) {
				case ErrUnauthorized:
					handlers.HandleResponseError(w, err.Error(), http.StatusUnauthorized)
					return
				default:
					handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			ctx := WithSession(r.Context(), &session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
