package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
	"github.com/Pashakrut94/SwiftChat/chat"
	"github.com/Pashakrut94/SwiftChat/users"
)

var (
	 pgUser = flag.String("pg_user", "Pasha","PostgreSQL name")
	 pgPwd = flag.String("pg_pwd", "pwd0123456789","PostgreSQL password")
	 pgHost = flag.String("pg_host","localhost","PostgreSQL host")	 
	 pgPort = flag.String("pg_port", "54320","PostgreSQL port")
	 pgDBname = flag.String("pg_dbname","mydb","PostgreSQL name of DB")
)

func main() {
	flag.Parse()
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",*pgUser, *pgPwd, *pgHost, *pgPort, *pgDBname)

	db, err := sql.Open("postgres", connectionString)
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
