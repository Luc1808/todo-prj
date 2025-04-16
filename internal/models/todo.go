package models

import "time"

type Todo struct {
	Title       string
	Description string
	Complete    bool
	Priority    string
	Category    string
	CreatedAt   time.Time
	DueAt       time.Time
	UserID      uint
}
