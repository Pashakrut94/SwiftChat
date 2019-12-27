package users

import (
	"net/http"

	_ "github.com/lib/pq"

	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/gorilla/mux"
)

func CreateUser(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handlers.HandleResponseError(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var user User
		if err := json.Unmarshal(body, &user); err != nil {
			handlers.HandleResponseError(w, "Error unmarshaling request body", http.StatusInternalServerError)
			return
		}
		if err := repo.Create(&user); err != nil {
			handlers.HandleResponseError(w, "Error creating new user", http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, user, pretty)
	}
}

func GetUser(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["UserID"])
		if err != nil {
			handlers.HandleResponseError(w, "Incorrect enter of UserID", http.StatusBadRequest)
			return
		}
		user, err := repo.Get(userID)
		if err != nil {
			handlers.HandleResponseError(w, "Error getting user by ID", http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, user, pretty)
	}
}

func ListUsers(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allUsers, err := repo.List()
		if err != nil {
			handlers.HandleResponseError(w, "Error listing of users", http.StatusInternalServerError)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, allUsers, pretty)
	}
}
