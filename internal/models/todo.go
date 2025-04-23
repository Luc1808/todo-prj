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
	ID          uint
	Title       string     `json:"title"`
	Description string     `json:"description"`
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
	query := `INSERT INTO todo (title, description, complete, priority, category, createdAt, dueDate, userID) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.DB.Exec(query, t.Title, t.Description, t.Complete, t.Priority, t.Category, time.Now(), t.DueAt, t.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todo) GetAllTodos() ([]Todo, error) {
	query := `SELECT id, title, description, complete, priority, category, createdat, duedate, userID FROM todo WHERE userid = $1`
	rows, err := db.DB.Query(query, t.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Complete, &todo.Priority, &todo.Category, &todo.CreatedAt, &todo.DueAt, &todo.UserID)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func GetCompletedTodos(userID uint) ([]Todo, error) {
	query := `SELECT id, title, description, complete, priority, category, createdat, duedate FROM todo WHERE userid = $1 AND complete = true`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Complete, &todo.Priority, &todo.Category, &todo.CreatedAt, &todo.DueAt)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

func (t *Todo) GetTodoByID(id uint) error {
	query := `SELECT title, description, complete, priority, category, createdat, duedate, userID FROM todo WHERE id = $1`
	row := db.DB.QueryRow(query, id)

	err := row.Scan(&t.Title, &t.Description, &t.Complete, &t.Priority, &t.Category, &t.CreatedAt, &t.DueAt, &t.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todo) UpdateTodo(id uint, userID uint) error {
	query := `UPDATE todo SET title = $1, description = $2, complete = $3, priority = $4, category = $5, duedate = $6 WHERE id = $7 AND userID = $8`
	_, err := db.DB.Exec(query, t.Title, t.Description, t.Complete, t.Priority, t.Category, t.DueAt, id, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todo) DeleteTodo(id uint, userID uint) error {
	query := `DELETE FROM todo WHERE id = $1 AND userID = $2`
	_, err := db.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	return nil
}
