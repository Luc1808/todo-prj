package models

import (
	"time"

	"github.com/Luc1808/todo-prj/internal/db"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	ID        uint
	Token     string
	UserID    uint
	CreatedAt time.Time
	ExpiresAt time.Time
}

func StoreRefreshToken(userID uint, token string, expiresAt time.Time) error {
	query := `INSERT INTO token (token, userID, createdAt, expiresAt) VALUES ($1, $2, $3, $4)`

	_, err := db.DB.Exec(query, token, userID, time.Now(), expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func FindRefreshToken(token string, userID uint) (*RefreshToken, error) {
	query := `SELECT * FROM token WHERE token = $1 AND userID = $2 AND expiresAt > $3`

	var fetchedToken RefreshToken
	err := db.DB.QueryRow(query, token, userID, time.Now()).Scan(
		&fetchedToken.ID,
		&fetchedToken.Token,
		&fetchedToken.UserID,
		&fetchedToken.CreatedAt,
		&fetchedToken.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return &fetchedToken, nil
}

func DeleteRefreshToken(token string) error {
	query := `DELETE FROM token WHERE token = $1`
	_, err := db.DB.Exec(query, token)
	if err != nil {
		return err
	}

	return nil
}
