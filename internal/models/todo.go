package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/Luc1808/todo-prj/internal/db"
)

type Priorities string
type Categories string

const (
	High   Priorities = "high"
	Medium Priorities = "medium"
	Low    Priorities = "low"
)

const (
	Health  Categories = "health"
	SelfDev Categories = "self development"
	Finance Categories = "finance"
	Social  Categories = "social"
)

type Todo struct {
	ID          uint
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Complete    bool       `json:"complete"`
	Priority    Priorities `json:"priority"`
	Category    Categories `json:"category"`
	CreatedAt   time.Time  `json:"createdAt"`
	DueDate     time.Time  `json:"dueDate"`
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
	_, err := db.DB.Exec(query, t.Title, t.Description, t.Complete, t.Priority, t.Category, time.Now(), t.DueDate, t.UserID)
	if err != nil {
		return err
	}

	return nil
}

// func (t *Todo) GetAllTodos() ([]Todo, error) {
// 	query := `SELECT id, title, description, complete, priority, category, createdat, duedate, userID FROM todo WHERE userid = $1`
// 	rows, err := db.DB.Query(query, t.UserID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var todos []Todo

// 	for rows.Next() {
// 		var todo Todo
// 		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Complete, &todo.Priority, &todo.Category, &todo.CreatedAt, &todo.DueDate, &todo.UserID)
// 		if err != nil {
// 			return nil, err
// 		}

// 		todos = append(todos, todo)
// 	}

// 	return todos, nil
// }

func GetTodosWithPagination(
	userID uint,
	limit int,
	offset int,
	sortBy string,
	sortOrder string,
	filterPriority Priorities,
	filterCategory Categories,
	filterComplete *bool,
	search string,
) ([]Todo, error) {
	var conditions []string
	var args []any
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("userid = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if filterPriority != "" {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, filterPriority)
		argIndex++
	}

	if filterCategory != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, filterCategory)
		argIndex++
	}

	if filterComplete != nil {
		conditions = append(conditions, fmt.Sprintf("complete = $%d", argIndex))
		args = append(args, *filterComplete)
		argIndex++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d OR description ILIKE $%d", argIndex, argIndex))
		args = append(args, "%"+search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var orderBy string
	if sortBy == "priority" {
		orderBy = fmt.Sprintf(
			`CASE priority
			WHEN '%s' THEN 1 
			WHEN '%s' THEN 2 
			WHEN '%s' THEN 3
			ELSE 4
			END %s`, High, Medium, Low, sortOrder)
	} else {
		orderBy = fmt.Sprintf("%s %s", sortBy, sortOrder)
	}

	query := fmt.Sprintf(
		`SELECT id, title, description, complete, priority, category, createdat, duedate 
		FROM todo 
		%s
		ORDER BY %s 
		LIMIT $%d OFFSET $%d 
		`, whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var todo Todo
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Complete, &todo.Priority, &todo.Category, &todo.CreatedAt, &todo.DueDate)
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

	err := row.Scan(&t.Title, &t.Description, &t.Complete, &t.Priority, &t.Category, &t.CreatedAt, &t.DueDate, &t.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todo) UpdateTodo(id uint, userID uint) error {
	query := `UPDATE todo SET title = $1, description = $2, complete = $3, priority = $4, category = $5, duedate = $6 WHERE id = $7 AND userID = $8`
	_, err := db.DB.Exec(query, t.Title, t.Description, t.Complete, t.Priority, t.Category, t.DueDate, id, userID)
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

// func GetTodoCount(userID uint) (int, error) {
// 	query := `SELECT COUNT(*) FROM todo WHERE userID = $1`
// 	var total int
// 	err := db.DB.QueryRow(query, userID).Scan(&total)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return total, nil
// }

func GetFilteredTodoCount(userID uint, priority Priorities, category Categories, complete *bool, search string) (int, error) {
	var conditions []string
	var args []any
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("userid = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if priority != "" {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, priority)
		argIndex++
	}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, category)
		argIndex++
	}

	if complete != nil {
		conditions = append(conditions, fmt.Sprintf("complete = $%d", argIndex))
		args = append(args, *complete)
		argIndex++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d OR description ILIKE $%d", argIndex, argIndex))
		args = append(args, "%"+search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM todo %s`, whereClause)

	var total int
	err := db.DB.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
