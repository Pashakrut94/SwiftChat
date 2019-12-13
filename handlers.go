package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

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
		data, err := FormatResp(user, pretty)
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
		data, err := FormatResp(allUsers, pretty)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func CreateRoom(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
			return
		}
		var room Room
		if err := json.Unmarshal(body, &room); err != nil {
			http.Error(w, "Error unmarshaling request body",
				http.StatusInternalServerError)
			return
		}
		if err := repo.Create(&room); err != nil {
			http.Error(w, "HTTP 400 Bad Request",
				http.StatusBadRequest)
			return
		}
		w.Write([]byte(strconv.Itoa(room.ID)))
	}
}

func ListRooms(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pretty := r.URL.Query()["pretty"]
		rooms, err := repo.List()
		if err != nil {
			http.Error(w, "Error listing of rooms",
				http.StatusInternalServerError)
			return
		}
		data, err := FormatResp(rooms, pretty)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func GetRoom(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID, err := strconv.Atoi(vars["RoomID"])
		if err != nil {
			http.Error(w, "Incorrect enter of RoomID",
				http.StatusBadRequest)
			return
		}
		room, err := repo.Get(roomID)
		if err != nil {
			http.Error(w, "Error getting room by ID",
				http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		data, err := FormatResp(room, pretty)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

type CreateMessageRequest struct {
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}

func CreateMessage(repo MsgRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID, err := strconv.Atoi(vars["RoomID"])
		if err != nil {
			http.Error(w, "Incorrect enter of RoomID",
				http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
			return
		}
		var req CreateMessageRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Error unmarshaling request body",
				http.StatusInternalServerError)
			return
		}
		msg := Message{Text: req.Text, UserID: req.UserID, RoomID: roomID}
		if err := repo.Create(&msg); err != nil {
			http.Error(w, "HTTP 400 Bad Request",
				http.StatusBadRequest)
			return
		}
		w.Write([]byte(strconv.Itoa(msg.ID)))
	}
}

func ListMessages(repo MsgRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID, err := strconv.Atoi(vars["RoomID"])
		if err != nil {
			http.Error(w, "Incorrect enter of RoomID",
				http.StatusBadRequest)
			return
		}
		msg, err := repo.ListByRoomID(roomID)
		if err != nil {
			http.Error(w, "Error getting messages by roomID",
				http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		data, err := FormatResp(msg, pretty)
		if err != nil {
			http.Error(w, "Error converting results to json",
				http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}
}

func FormatResp(payload interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(payload, "", " ")
	}
	return json.Marshal(payload)
}
