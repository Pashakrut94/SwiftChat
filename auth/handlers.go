package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Pashakrut94/SwiftChat/users"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func SignUp(userRepo users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(sessionName)
		if err == nil {
			http.Error(w, "This user already has a session",
				http.StatusInternalServerError)
			return
		}
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Wrong structure of the body in decoding",
				http.StatusBadRequest)
			return
		}
		if len(req.Password) < 8 {
			http.Error(w, "Password can't be shorter then 8 symbols",
				http.StatusBadRequest)
			return
		}
		if len(req.Phone) < 12 {
			http.Error(w, "Incorrect phone number",
				http.StatusBadRequest)
			return
		}
		if len(req.Username) < 2 {
			http.Error(w, "Username is too short",
				http.StatusBadRequest)
			return
		}
		_, err = userRepo.GetByPhone(req.Phone)
		if err == nil {
			http.Error(w, "Such phone number is already in use", 500)
			return
		}
		if err != sql.ErrNoRows {
			http.Error(w, "Error getting user by phone", 500)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}
		user := users.User{Name: req.Username, Phone: req.Phone, Password: string(hashedPassword)}
		if err := userRepo.Create(&user); err != nil {
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}
		if err := SetNewSession(sessRepo, user.ID, w); err != nil {
			http.Error(w, "Error setting session",
				http.StatusInternalServerError)
			return
		}
	}
}

type SignInRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func SignIn(userRepo users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(sessionName)
		if err == nil {
			http.Error(w, "This user already has a session",
				http.StatusInternalServerError)
			return
		}
		var req SignInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Wrong structure of the body in decoding",
				http.StatusBadRequest)
			return
		}
		if len(req.Phone) == 0 {
			http.Error(w, "Put the phone number",
				http.StatusBadRequest)
			return
		}
		if len(req.Password) == 0 {
			http.Error(w, "Put the password",
				http.StatusBadRequest)
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
			http.Error(w, "Incorrect password",
				http.StatusInternalServerError)
			return
		}
		if err := SetNewSession(sessRepo, user.ID, w); err != nil {
			http.Error(w, "Error setting session",
				http.StatusInternalServerError)
			return
		}
	}
}

func Logout(userRepo users.UserRepo, sessRepo SessionRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionName)
		if err != nil {
			http.Error(w, "Error getting session",
				http.StatusUnauthorized)
			return
		}
		if err := sessRepo.Delete(c.Value); err != nil {
			http.Error(w, "Error deleting session",
				http.StatusInternalServerError)
			return
		}
		c = &http.Cookie{
			Name:     sessionName,
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}
}
