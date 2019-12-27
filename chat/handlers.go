package chat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Pashakrut94/SwiftChat/auth"
	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func CreateRoom(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pretty := r.URL.Query()["pretty"]
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handlers.HandleResponseError(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var room Room
		if err := json.Unmarshal(body, &room); err != nil {
			handlers.HandleResponseError(w, "Error unmarshaling request body", http.StatusInternalServerError)
			return
		}
		if err := repo.Create(&room); err != nil {
			handlers.HandleResponseError(w, "HTTP 400 Bad Request", http.StatusBadRequest)
			return
		}
		handlers.HandleResponse(w, room, pretty)
	}
}

func ListRooms(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pretty := r.URL.Query()["pretty"]
		rooms, err := repo.List()
		if err != nil {
			handlers.HandleResponseError(w, "Error listing of rooms", http.StatusNotFound)
			return
		}
		handlers.HandleResponse(w, rooms, pretty)
	}
}

func GetRoom(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID, err := strconv.Atoi(vars["RoomID"])
		if err != nil {
			handlers.HandleResponseError(w, "Incorrect enter of RoomID", http.StatusBadRequest)
			return
		}
		room, err := repo.Get(roomID)
		if err != nil {
			handlers.HandleResponseError(w, "Error getting room by ID", http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, room, pretty)
	}
}

type CreateMessageRequest struct {
	Text string `json:"text"`
}

func CreateMessage(repo MsgRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := auth.SessionValue(ctx)
		vars := mux.Vars(r)
		roomID, err := strconv.Atoi(vars["RoomID"])
		if err != nil {
			handlers.HandleResponseError(w, "Incorrect enter of RoomID", http.StatusBadRequest)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handlers.HandleResponseError(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		var req CreateMessageRequest
		if err := json.Unmarshal(body, &req); err != nil {
			handlers.HandleResponseError(w, "Error unmarshaling request body", http.StatusInternalServerError)
			return
		}
		msg := Message{Text: req.Text, UserID: session.UserID, RoomID: roomID}
		if err := repo.Create(&msg); err != nil {
			handlers.HandleResponseError(w, "Error creating message", http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, msg, pretty)
	}
}

func ListMessages(repo MsgRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roomID, err := strconv.Atoi(vars["RoomID"])
		if err != nil {
			handlers.HandleResponseError(w, "Incorrect enter of RoomID", http.StatusBadRequest)
			return
		}
		msg, err := repo.ListByRoomID(roomID)
		if err != nil {
			handlers.HandleResponseError(w, "Error getting messages by roomID", http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, msg, pretty)
	}
}
