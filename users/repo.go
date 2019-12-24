package users

import (
	"database/sql"
	"fmt"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) get(q string, args ...interface{}) (*User, error) {
	var u User
	if err := repo.db.QueryRow(q, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Phone,
		&u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

func (repo *UserRepo) Get(id int) (*User, error) {
	return repo.get("select * from users where id = $1", id)
}

func (repo *UserRepo) GetByPhone(phone string) (*User, error) {
	return repo.get("select * from users where phone = $1", phone)
}

func (repo *UserRepo) List() ([]User, error) {
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
	return users, nil
}

func (repo *UserRepo) Create(user *User) error {
	err := repo.db.QueryRow("insert into users (name, phone, password) values ($1,$2,$3) returning id", user.Name, user.Phone, user.Password).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
