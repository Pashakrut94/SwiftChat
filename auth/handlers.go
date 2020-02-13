package auth

import (
	"database/sql"
	"encoding/json"
	"text/template"

	"net/http"
	"time"

	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/Pashakrut94/SwiftChat/users"
	"github.com/pkg/errors"
)

func SignUpTemplate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(sessionName)
		if err == nil {
			http.Redirect(w, r, "/api/hellopage", http.StatusTemporaryRedirect)
		}
		t, err := template.ParseFiles("./auth/templates/signup.gtpl")
		if err != nil {
			handlers.HandleResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, nil)
	}
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func FormParseSignUp(r *http.Request) (*SignUpRequest, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	name := r.FormValue("username")
	password := r.FormValue("password")
	phone := r.FormValue("phone")
	req := SignUpRequest{Username: name, Password: password, Phone: phone}
	return &req, nil
}

func CheckCookies(userRepo *users.UserRepo, w http.ResponseWriter, r *http.Request, phone string) error {
	_, err := r.Cookie(sessionName)
	if err == nil {
		http.Redirect(w, r, "/api/hellopage", http.StatusTemporaryRedirect)
		return nil
	} else {
		_, err := userRepo.GetByPhone(phone)
		switch err {
		case nil:
			return errors.New("user already exists")
		case sql.ErrNoRows:
			return nil
		default:
			return err
		}
	}
}

func SignUp(userRepo *users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := FormParseSignUp(r)
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error parsing template").Error(), http.StatusInternalServerError)
			return
		}

		// if err := CheckCookies(userRepo, w, r, req.Phone); err != nil {
		// 	handlers.HandleResponseError(w, err.Error(), http.StatusBadGateway) // another http Error
		// 	return
		// }

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
		if err == nil { //if session exists redirect to main page on form method "get" do it
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
