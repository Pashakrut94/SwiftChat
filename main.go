package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/Pashakrut94/SwiftChat/chat"
	"github.com/Pashakrut94/SwiftChat/users"
)

func main() {

	connStr := "postgresql://Pasha:pwd0123456789@localhost:54320/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := users.NewUserRepo(db)
	roomRepo := chat.NewRoomRepo(db)
	msgRepo := chat.NewMsgRepo(db)

	router := mux.NewRouter()

	router.HandleFunc("/api/rooms", chat.CreateRoom(*roomRepo)).Methods("POST")
	router.HandleFunc("/api/rooms", chat.ListRooms(*roomRepo)).Methods("GET")
	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}", chat.GetRoom(*roomRepo)).Methods("GET")

	router.HandleFunc("/api/users", users.CreateUsers(*userRepo)).Methods("POST")
	router.HandleFunc("/api/users", users.ListUsers(*userRepo)).Methods("GET")
	router.HandleFunc("/api/users/{UserID:[0-9]+}", users.GetUser(*userRepo)).Methods("GET")

	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}/messages", chat.CreateMessage(*msgRepo)).Methods("POST")
	router.HandleFunc("/api/rooms/{RoomID:[0-9]+}/messages", chat.ListMessages(*msgRepo)).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Server starts at :8080")
	http.ListenAndServe(":8080", nil)
}
