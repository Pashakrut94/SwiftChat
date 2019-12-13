package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type Message struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
	RoomID int    `json:"room_id"`
}

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {

	connStr := "postgresql://Pasha:pwd0123456789@localhost:54320/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := UserRepo{db: db}
	roomRepo := RoomRepo{db: db}
	msgRepo := MsgRepo{db: db}

	router := mux.NewRouter()

	router.HandleFunc("/api/rooms", CreateRoom(roomRepo)).Methods("POST")
	router.HandleFunc("/api/rooms", ListRooms(roomRepo)).Methods("GET")
	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}", GetRoom(roomRepo)).Methods("GET")

	router.HandleFunc("/api/users", CreateUsers(userRepo)).Methods("POST")
	router.HandleFunc("/api/users", ListUsers(userRepo)).Methods("GET")
	router.HandleFunc("/api/users/{UserID:[0-9]+}", GetUser(userRepo)).Methods("GET")

	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}/messages", CreateMessage(msgRepo)).Methods("POST")
	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}/messages", ListMessages(msgRepo)).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Server starts at :8080")
	http.ListenAndServe(":8080", nil)
}
