package models

import (
	"errors"

	"github.com/Luc1808/todo-prj/internal/db"
	"github.com/Luc1808/todo-prj/internal/utils"
)

type User struct {
	ID       uint
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Save() error {
	query := `INSERT INTO users (email, password) VALUES ($1, $2)`

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	// Maybe scan back the user ID later
	_, err = db.DB.Exec(query, u.Email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) VerifyCredentials() error {
	query := `SELECT id, password FROM users WHERE email = $1`

	var dbPassword string
	err := db.DB.QueryRow(query, u.Email).Scan(&u.ID, &dbPassword)
	if err != nil {
		return err
	}

	isPassValid := utils.CheckPasswordHash(u.Password, dbPassword)
	if !isPassValid {
		return errors.New("invalid credentials")
	}

	return nil
}

func GetAllUsers() ([]User, error) {
	query := `SELECT * FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Email, &user.Password)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
