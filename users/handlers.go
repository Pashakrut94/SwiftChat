package users

import (
	"net/http"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"strconv"

	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/gorilla/mux"
)

func GetUser(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["UserID"])
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "incorrect enter of userID").Error(), http.StatusBadRequest)
			return
		}
		user, err := HandleGetUSer(repo, userID)
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, err.Error()).Error(), http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, user, pretty)
	}
}

func ListUsers(repo UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allUsers, err := HandleListUSer(repo)
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, err.Error()).Error(), http.StatusInternalServerError)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, allUsers, pretty)
	}
}
