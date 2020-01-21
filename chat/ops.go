package chat

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

func validateCreateRoom(repo RoomRepo, room Room) error {
	_, err := repo.Get(room.ID)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "error getting room from DB")
	}
	return nil
}

func HandleCreateRoom(repo RoomRepo, room Room) (Room, error) {
	if err := validateCreateRoom(repo, room); err != nil {
		return Room{}, err
	}
	if err := repo.Create(&room); err != nil {
		fmt.Println(err)
		return Room{}, err
	}
	return room, nil
}

func HandleListRooms(repo RoomRepo) ([]Room, error) {
	rooms, err := repo.List()
	if err == sql.ErrNoRows {
		return []Room{}, nil
	}
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func HandleGetRoom(repo RoomRepo, roomID int) (Room, error) {
	room, err := repo.Get(roomID)
	if err == sql.ErrNoRows {
		return Room{}, ErrNotFound
	}
	if err != nil {
		return Room{}, err
	}
	return *room, nil
}

func HandleCreateMessage(repo MsgRepo, text string, userID, roomID int) (Message, error) {
	msg := Message{Text: text, UserID: userID, RoomID: roomID}
	if err := repo.Create(&msg); err != nil {
		return Message{}, errors.New("error creating message")
	}
	return msg, nil
}

func HandleListMessages(repo MsgRepo, roomID int) ([]Message, error) {
	msgs, err := repo.ListByRoomID(roomID)
	if err != nil {
		return nil, errors.Wrap(err, "error getting messages from DB")
	}
	return msgs, nil
}
