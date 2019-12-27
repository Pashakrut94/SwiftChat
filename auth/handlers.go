package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/Pashakrut94/SwiftChat/users"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func ConditionsToSignUp(w http.ResponseWriter, r *http.Request, req SignUpRequest) bool {
	if len(req.Username) < 2 {
		handlers.HandleResponseError(w, "Username is too short", http.StatusBadRequest)
		return false
	}
	if len(req.Password) < 8 {
		handlers.HandleResponseError(w, "Password can't be shorter then 8 symbols", http.StatusBadRequest)
		return false
	}
	if len(req.Phone) < 12 {
		handlers.HandleResponseError(w, "Incorrect phone number", http.StatusBadRequest)
		return false
	}
	return true
}

func SignUp(userRepo users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(sessionName)
		if err == nil {
			handlers.HandleResponseError(w, "This user already has a session", http.StatusInternalServerError)
			return
		}
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.HandleResponseError(w, "Wrong structure of the body in decoding", http.StatusBadRequest)
			return
		}
		canSignUp := ConditionsToSignUp(w, r, req)
		if !canSignUp {
			return
		}
		_, err = userRepo.GetByPhone(req.Phone)
		if err == nil {
			handlers.HandleResponseError(w, "Such phone number is already in use", http.StatusInternalServerError)
			return
		}
		if err != sql.ErrNoRows {
			handlers.HandleResponseError(w, "There is no such phone number in DB", http.StatusInternalServerError)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			handlers.HandleResponseError(w, "Server error, unable to create your account", http.StatusInternalServerError)
			return
		}
		user := users.User{Name: req.Username, Phone: req.Phone, Password: string(hashedPassword)}
		if err := userRepo.Create(&user); err != nil {
			handlers.HandleResponseError(w, "Server error, unable to create your account", http.StatusInternalServerError)
			return
		}
		SetNewSession(w, sessRepo, user.ID)
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, user, pretty)
	}
}

type SignInRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func ConditionsToSignIn(w http.ResponseWriter, r *http.Request, req SignInRequest) bool {
	if len(req.Phone) < 12 {
		handlers.HandleResponseError(w, "Incorrect input of phone number", http.StatusBadRequest)
		return false
	}
	if len(req.Password) < 8 {
		handlers.HandleResponseError(w, "Password can't be shorter then 8 symbols", http.StatusBadRequest)
		return false
	}
	return true
}

func SignIn(userRepo users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(sessionName)
		if err == nil {
			handlers.HandleResponseError(w, "This user already has a session", http.StatusInternalServerError)
			return
		}
		var req SignInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.HandleResponseError(w, "Wrong structure of the body in decoding", http.StatusBadRequest)
			return
		}
		canSignIn := ConditionsToSignIn(w, r, req)
		if !canSignIn {
			return
		}
		user, err := userRepo.GetByPhone(req.Phone)
		if err != nil {
			http.Error(w, "Incorrect phone number",
				http.StatusBadRequest)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			handlers.HandleResponseError(w, "Incorrect password", http.StatusInternalServerError)
			return
		}
		SetNewSession(w, sessRepo, user.ID)
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, user, pretty)
	}
}

func Logout(userRepo users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionName)
		if err != nil {
			handlers.HandleResponseError(w, "Error getting session", http.StatusUnauthorized)
			return
		}
		if err := sessRepo.Delete(c.Value); err != nil {
			handlers.HandleResponseError(w, "Error deleting session", http.StatusInternalServerError)
			return
		}
		c = &http.Cookie{
			Name:     sessionName,
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, c, pretty)
	}
}
