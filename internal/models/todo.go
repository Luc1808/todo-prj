package models

import (
	"time"

	"github.com/Luc1808/todo-prj/internal/db"
)

type Priorities string
type Categories string

const (
	High   Priorities = "High"
	Medium Priorities = "Medium"
	Low    Priorities = "Low"
)

const (
	Health  Categories = "Health"
	SelfDev Categories = "Self Development"
	Finance Categories = "Finance"
	Social  Categories = "Social"
)

type Todo struct {
	ID          string
	Title       string     `json:"title"`
	Description string     `json:"descrition"`
	Complete    bool       `json:"complete"`
	Priority    Priorities `json:"priority"`
	Category    Categories `json:"category"`
	CreatedAt   time.Time  `json:"createdAt"`
	DueAt       time.Time  `json:"dueAt"`
	UserID      uint       `json:"userID"`
}

func (p Priorities) IsValid() bool {
	switch p {
	case High, Medium, Low:
		return true
	}
	return false
}

func (c Categories) IsValid() bool {
	switch c {
	case Health, SelfDev, Finance, Social:
		return true
	}
	return false
}

func (t *Todo) Save() error {
	query := `INSERT INTO todo (title, description, complete, priority, category, createdat, duedate, userid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.DB.Exec(query, t.Title, t.Description, t.Complete, t.Priority, t.Category, time.Now(), t.DueAt, t.UserID)
	if err != nil {
		return err
	}

	return nil
}
