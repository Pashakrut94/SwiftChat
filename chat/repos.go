package chat

import (
	"database/sql"
	"fmt"
)

type MsgRepo struct {
	db *sql.DB
}

func NewMsgRepo(db *sql.DB) *MsgRepo {
	return &MsgRepo{db: db}
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

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{db: db}
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
