package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	RoomID int    `json:"room_id"`
	Sender User   `json:"sender"`
	Text   string `json:"text"`
}

type Room struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Users    []User    `json:"users"`
	Messages []Message `json:"messages"`
}

var rooms = make(map[int]Room)

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	var room Room
	if err := json.Unmarshal(body, &room); err != nil {
		http.Error(w, "Error unmarshaling request body",
			http.StatusInternalServerError)
		return
	}
	rooms[room.ID] = room
}
func ListRooms(w http.ResponseWriter, r *http.Request) {
	_, pretty := r.URL.Query()["pretty"]
	data, err := FormatResp(rooms, pretty)
	if err != nil {
		http.Error(w, "Error converting results to json",
			http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func GetRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID := vars["RoomID"]
	id, err := strconv.Atoi(roomID)
	if err != nil {
		http.Error(w, "Incorrect enter of RoomID",
			http.StatusBadRequest)
		return
	}
	_, pretty := r.URL.Query()["pretty"]
	data, err := FormatResp(rooms[id], pretty)
	if err != nil {
		http.Error(w, "Error converting results to json",
			http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body",
			http.StatusInternalServerError)
		return
	}
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Error unmarshaling request body",
			http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	room := rooms[msg.RoomID]
	room.Messages = append(room.Messages, msg)
	rooms[msg.RoomID] = room
}
func ListMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	number := vars["RoomID"]
	id, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, "Incorrect enter of RoomID",
			http.StatusBadRequest)
		return
	}
	_, pretty := r.URL.Query()["pretty"]
	data, err := FormatResp(rooms[id].Messages, pretty)
	if err != nil {
		http.Error(w, "Error converting results to json",
			http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
func FormatResp(payload interface{}, pretty bool) ([]byte, error) {
	if pretty {
		return json.MarshalIndent(payload, "", " ")
	}
	return json.Marshal(payload)
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/rooms", CreateRoom).Methods("POST")
	router.HandleFunc("/api/rooms", ListRooms).Methods("GET")

	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}", GetRoom).Methods("GET")

	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}/messages", CreateMessage).Methods("POST")
	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}/messages", ListMessages).Methods("GET")

	// Note: Expected message: { text: "sample text", sender: 14 }.
	// TODO(pasha): CreateMessageRequest struct to unmarshal new message json format.

	http.Handle("/", router)
	fmt.Println("Server starts at :8080")
	http.ListenAndServe(":8080", nil)
}
