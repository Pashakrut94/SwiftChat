package main

import (
	"database/sql"
	"fmt"
)

type MsgRepo struct {
	db *sql.DB
}

func (repo *MsgRepo) Create(msg *Message) error {
	q := "insert into messages (text, user_id, room_id) values ($1, $2, $3) returning id"
	err := repo.db.QueryRow(q, msg.Text, msg.UserID, msg.RoomID).Scan(&msg.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MsgRepo) ListByRoomID(roomID int) ([]Message, error) {
	q := "select id, text, user_id, room_id from messages where room_id = $1"
	rows, err := repo.db.Query(q, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var msgs []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(&msg.ID, &msg.Text, &msg.UserID, &msg.RoomID)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if msg.RoomID == roomID {
			msgs = append(msgs, msg)
		}
	}
	return msgs, nil
}

type RoomRepo struct {
	db *sql.DB
}

func (repo *RoomRepo) Create(room *Room) error {
	q := "insert into rooms (name) values ($1) returning id"
	err := repo.db.QueryRow(q, room.Name).Scan(&room.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *RoomRepo) List() ([]Room, error) {
	q := "select id , name from rooms"
	rows, err := repo.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rooms []Room
	for rows.Next() {
		var room Room
		err := rows.Scan(&room.ID, &room.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (repo *RoomRepo) Get(id int) (*Room, error) {
	row := repo.db.QueryRow("select id , name from rooms where id = $1", id)
	var room Room
	err := row.Scan(&room.ID, &room.Name)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

type UserRepo struct {
	db *sql.DB
}

func (repo *UserRepo) Get(id int) (*User, error) {
	row := repo.db.QueryRow("select * from users where id = $1", id)
	user := User{}
	err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) List() (*[]User, error) {
	rows, err := repo.db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []User{}
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Phone, &user.Password)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, user)
	}
	return &users, nil
}

func (repo *UserRepo) Create(user *User) error {
	err := repo.db.QueryRow("insert into users (name, phone, password) values ($1,$2,$3) returning id", user.Name, user.Phone, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
