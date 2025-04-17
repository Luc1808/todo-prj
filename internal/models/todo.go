package models

import "time"

type Todo struct {
	Title       string    `json:"title"`
	Description string    `json:"descrition"`
	Complete    bool      `json:"complete"`
	Priority    string    `json:"priority"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"createdAt"`
	DueAt       time.Time `json:"dueAt"`
	UserID      uint      `json:"userID"`
}
