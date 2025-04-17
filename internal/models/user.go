package models

import (
	"github.com/Luc1808/todo-prj/internal/db"
	"github.com/Luc1808/todo-prj/internal/utils"
)

type User struct {
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
