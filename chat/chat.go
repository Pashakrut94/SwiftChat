package chat

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
