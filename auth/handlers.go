package auth

import (
	"encoding/json"

	"net/http"
	"time"

	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/Pashakrut94/SwiftChat/users"
	"github.com/pkg/errors"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func SignUp(userRepo *users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(sessionName)
		if err == nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "this user already has a session").Error(), http.StatusInternalServerError)
			return
		}
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error parsing signup request").Error(), http.StatusBadRequest)
			return
		}
		user, err := HandleSignUp(userRepo, req.Username, req.Password, req.Phone)
		if err != nil {
			switch errors.Cause(err) {
			case ErrValidationError, ErrUserAlreadyExists:
				handlers.HandleResponseError(w, err.Error(), http.StatusBadRequest)
				return
			default:
				handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		sid, expires, err := SetNewSession(sessRepo, user.ID)
		if err != nil {
			handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: sessionName, Value: sid, HttpOnly: true, Expires: expires})
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, user, pretty)
	}
}

type SignInRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func SignIn(userRepo *users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(sessionName)
		if err == nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "this user already has a session").Error(), http.StatusInternalServerError)
			return
		}
		var req SignInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "wrong structure of the body in decoding").Error(), http.StatusBadRequest)
			return
		}
		user, err := HandleSignIn(userRepo, req.Phone, req.Password)
		if err != nil {
			switch errors.Cause(err) {
			case ErrValidationError:
				handlers.HandleResponseError(w, err.Error(), http.StatusBadRequest)
				return
			default:
				handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		sid, expires, err := SetNewSession(sessRepo, user.ID)
		if err != nil {
			handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: sessionName, Value: sid, HttpOnly: true, Expires: expires})
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, user, pretty)
	}
}

func Logout(userRepo users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionName)
		if err != nil {
			handlers.HandleResponseError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err := sessRepo.Delete(c.Value); err != nil {
			handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c = &http.Cookie{
			Name:     sessionName,
			Value:    "",
			Expires:  time.Unix(0, 0).UTC(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)

		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, nil, pretty)
	}
}
