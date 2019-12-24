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

func CreateUsers(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
			return
		}
		var user User
		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "Error unmarshaling request body",
				http.StatusInternalServerError)
			return
		}
		if err := repo.Create(&user); err != nil {
			http.Error(w, "HTTP 400 Bad Request",
				http.StatusBadRequest)
		}
		w.Write([]byte(strconv.Itoa(user.ID)))
	}
}

func GetUser(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["UserID"])
		if err != nil {
			http.Error(w, "Incorrect enter of UserID",
				http.StatusBadRequest)
			return
		}
		user, err := repo.Get(userID)
		if err != nil {
			http.Error(w, "Error getting user by ID",
				http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		data, err := handlers.FormatResp(user, pretty)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func ListUsers(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allUsers, err := repo.List()
		if err != nil {
			http.Error(w, "Error listing of users",
				http.StatusInternalServerError)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		data, err := handlers.FormatResp(allUsers, pretty)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}
