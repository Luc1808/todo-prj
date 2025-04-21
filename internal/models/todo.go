package models

import "time"

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
