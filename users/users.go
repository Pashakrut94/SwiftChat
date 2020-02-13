package users

import "database/sql"

type User struct {
	ID       int            `json:"id"`
	Name     string         `json:"name"`
	Password string         `json:"password"`
	Phone    string         `json:"phone"`
	Avatar   sql.NullString `json:"url"`
}
