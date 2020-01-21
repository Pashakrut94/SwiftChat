package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	"github.com/Pashakrut94/SwiftChat/auth"
	"github.com/Pashakrut94/SwiftChat/chat"
	"github.com/Pashakrut94/SwiftChat/users"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	pgUser   = flag.String("pg_user", "Pasha", "PostgreSQL name")
	pgPwd    = flag.String("pg_pwd", "pwd0123456789", "PostgreSQL password")
	pgHost   = flag.String("pg_host", "localhost", "PostgreSQL host")
	pgPort   = flag.String("pg_port", "54320", "PostgreSQL port")
	pgDBname = flag.String("pg_dbname", "mydb", "PostgreSQL name of DB")
)

func main() {
	flag.Parse()
	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", *pgUser, *pgPwd, *pgHost, *pgPort, *pgDBname)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	roomRepo := chat.NewRoomRepo(db)
	userRepo := users.NewUserRepo(db)
	msgRepo := chat.NewMsgRepo(db)
	sessRepo := auth.NewSessionRepo(db)
	authMiddleware := auth.RequireAuthentication(*sessRepo)

	router := mux.NewRouter()

	router.HandleFunc("/api/signup", auth.SignUp(userRepo, *sessRepo)).Methods("POST")
	router.HandleFunc("/api/signin", auth.SignIn(userRepo, *sessRepo)).Methods("POST")
	router.HandleFunc("/api/logout", auth.Logout(*userRepo, *sessRepo)).Methods("POST")

	router.Handle("/api/rooms", authMiddleware(chat.CreateRoom(*roomRepo))).Methods("POST")
	router.Handle("/api/rooms", authMiddleware(chat.ListRooms(*roomRepo))).Methods("GET")
	router.Handle("/api/rooms/{RoomID:[0-9]+}", authMiddleware(chat.GetRoom(*roomRepo))).Methods("GET")

	router.Handle("/api/users", authMiddleware(users.ListUsers(*userRepo))).Methods("GET")
	router.Handle("/api/users/{UserID:[0-9]+}", authMiddleware(users.GetUser(*userRepo))).Methods("GET")

	router.Handle("/api/rooms/{RoomID:[0-9]+}/messages", authMiddleware(chat.CreateMessage(*msgRepo))).Methods("POST")
	router.Handle("/api/rooms/{RoomID:[0-9]+}/messages", authMiddleware(chat.ListMessages(*msgRepo))).Methods("GET")

	http.Handle("/", router)
	fmt.Println("Server starts at :8080")
	http.ListenAndServe(":8080", nil)
}
