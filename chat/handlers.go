package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Pashakrut94/SwiftChat/auth"
	"github.com/Pashakrut94/SwiftChat/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

//уникальное имя, добавь кейс в ошибках для этого, добавь ошибку в этом пакете отдельную и обрабатывай через нее
func CreateRoom(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pretty := r.URL.Query()["pretty"]
		var room Room
		if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error parsing signup request").Error(), http.StatusBadRequest)
			return
		}
		room, err := HandleCreateRoom(repo, room)
		if err != nil {
			switch errors.Cause(err) {
			case ErrRoomAlreadyExists:
				handlers.HandleResponseError(w, errors.Wrap(err, "room already exists").Error(), http.StatusBadRequest)
				return
			default:
				handlers.HandleResponseError(w, errors.Wrap(err, "error creating room").Error(), http.StatusInternalServerError)
				return
			}
		}
		handlers.HandleResponse(w, room, pretty)
	}
}

func ListRooms(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pretty := r.URL.Query()["pretty"]
		rooms, err := HandleListRooms(repo)
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error listing of rooms").Error(), http.StatusInternalServerError)
			return
		}
		handlers.HandleResponse(w, rooms, pretty)
	}
}

func GetRoom(repo RoomRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		variables := mux.Vars(r)
		roomID, err := strconv.Atoi(variables["RoomID"])
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "incorrect enter of RoomID").Error(), http.StatusBadRequest)
			return
		}
		room, err := HandleGetRoom(repo, roomID)
		if err != nil {
			switch errors.Cause(err) {
			case ErrNotFound:
				handlers.HandleResponseError(w, errors.Wrap(err, err.Error()).Error(), http.StatusNotFound)
				return
			default:
				handlers.HandleResponseError(w, errors.Wrap(err, "error no rows").Error(), http.StatusInternalServerError)
				return
			}
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
		variables := mux.Vars(r)
		roomID, err := strconv.Atoi(variables["RoomID"])
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "incorrect enter of RoomID").Error(), http.StatusBadRequest)
			return
		}
		var req CreateMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error parsing signup request").Error(), http.StatusBadRequest)
			return
		}
		msg, err := HandleCreateMessage(repo, req.Text, session.UserID, roomID)
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, err.Error()).Error(), http.StatusInternalServerError)
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
			handlers.HandleResponseError(w, errors.Wrap(err, "incorrect enter of RoomID").Error(), http.StatusBadRequest)
			return
		}
		msgs, err := HandleListMessages(repo, roomID)
		if err != nil {
			handlers.HandleResponseError(w, errors.Wrap(err, "error getting messages").Error(), http.StatusBadRequest)
			return
		}
		_, pretty := r.URL.Query()["pretty"]
		handlers.HandleResponse(w, msgs, pretty)
	}
}
