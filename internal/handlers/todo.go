package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Luc1808/todo-prj/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type PaginatedTodosResponse struct {
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalTodos int           `json:"total_todos"`
	TotalPages int           `json:"total_pages"`
	SortBy     string        `json:"sort_by"`
	Order      string        `json:"order"`
	Todos      []models.Todo `json:"todos"`
}

// CreateTodoHandler godoc
// @Summary      Create a new todo
// @Description  Create a new todo task
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo  body     models.TodoRequest  true  "Todo object"
// @Success      201  {string}  string  "Todo successfully created"
// @Failure      400  {string}  string  "Invalid JSON"
// @Failure      401  {string}  string  "Unauthorized - missing token"
// @Failure      500  {string}  string  "Error trying to saving new To-do task"
// @Security BearerToken
// @Router       /tasks [post]
func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID for todo creation
	userID, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	// Get request body
	var todo models.Todo
	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !todo.Priority.IsValid() {
		http.Error(w, "Invalid priority value", http.StatusBadRequest)
		return
	}
	if !todo.Category.IsValid() {
		http.Error(w, "Invalid category value", http.StatusBadRequest)
		return
	}
	fmt.Println(userID)
	todo.UserID = userID

	err = todo.Save()
	if err != nil {
		http.Error(w, "Error trying to saving new To-do task", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Todo created successfully!")
}

// GetAllTodos is commented out because it is not used in the current implementation.

// func GetAllTodos(w http.ResponseWriter, r *http.Request) {
// 	authHeader := r.Header.Get("Authorization")

// 	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
// 		return
// 	}

// 	token := strings.TrimPrefix(authHeader, "Bearer ")
// 	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
// 		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
// 	})
// 	if err != nil || !parsedToken.Valid {
// 		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
// 		return
// 	}

// 	claims, ok := parsedToken.Claims.(jwt.MapClaims)
// 	if !ok {
// 		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
// 		return
// 	}

// 	userID := uint(claims["user_id"].(float64))

// 	var todo models.Todo
// 	todo.UserID = userID
// 	todos, err := todo.GetAllTodos()
// 	if err != nil {
// 		http.Error(w, "Error trying to get all todos", http.StatusBadRequest)
// 		return
// 	}

// 	w.Header().Set("Content-type", "application/json")
// 	json.NewEncoder(w).Encode(todos)
// }

// GetAllTodosWithPagination godoc
// @Summary      Get all todos with pagination
// @Description  Get all todos with pagination, sorting, and filtering
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        page  query     int  false  "Page number"
// @Param        limit  query     int  false  "Number of items per page"
// @Param        sort_by  query     string  false  "Sort by field (priority, category, createdat, duedate)"
// @Param        order  query     string  false  "Sort order (asc/desc)"
// @Param        priority  query     string  false  "Filter by priority (high, medium, low)"
// @Param        category  query     string  false  "Filter by category (social, self development, finance, health)"
// @Param        complete  query     bool  false  "Filter by completion status"
// @Param        search  query     string  false  "Search term"
// @Success      200  {object}  PaginatedTodosResponse
// @Failure      400  {string}  string  "Invalid parameters"
// @Failure      401  {string}  string  "Unauthorized - missing token"
// @Failure      500  {string}  string  "Failed to get the total number of to-dos"
// @Security BearerToken
// @Router       /tasks [get]
func GetAllTodosWithPagination(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	// Pagination, sorting, filtering stuff

	page, limit, offset := parsePaginationParams(r)
	sortBy, sortOrder := parseSortParams(r)
	priority, category, complete := parseFilterParams(r)
	search := r.URL.Query().Get("search")

	// totalTodos, err := models.GetTodoCount(userID)
	totalTodos, err := models.GetFilteredTodoCount(userID, priority, category, complete, search)
	if err != nil {
		http.Error(w, "Failed to get the total number of to-dos", http.StatusUnauthorized)
		return
	}

	todos, err := models.GetTodosWithPagination(userID, limit, offset, sortBy, sortOrder, priority, category, complete, search)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range todos {
		todos[i].UserID = userID
	}

	// Get total pages
	totalPages := int(math.Ceil(float64(totalTodos) / float64(limit)))

	// response := map[string]any{
	// 	"page":        page,
	// 	"limit":       limit,
	// 	"total_todos": totalTodos,
	// 	"total_pages": totalPages,
	// 	"sort_by":     sortBy,
	// 	"order":       sortOrder,
	// 	"todos":       todos,
	// }

	response := PaginatedTodosResponse{
		Page:       page,
		Limit:      limit,
		TotalTodos: totalTodos,
		TotalPages: totalPages,
		SortBy:     sortBy,
		Order:      sortOrder,
		Todos:      todos,
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTodoByID godoc
// @Summary      Get a todo by ID
// @Description  Get a todo task by its ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id  path     int  true  "Todo ID"
// @Success      200  {object}  models.Todo
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      401  {string}  string  "Unauthorized - missing token"
// @Failure      500  {string}  string  "Error trying to get todo by ID"
// @Security BearerToken
// @Router       /tasks/{id} [get]
func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	var todo models.Todo
	todo.UserID = userID

	todoIDStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	idUint64, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	todoID := uint(idUint64)
	todo.ID = uint(todoID)

	err = todo.GetTodoByID(todoID)
	if err != nil {
		http.Error(w, "Error trying to get todo by ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	var todo models.Todo
	todo.UserID = userID

	todoIDStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	idUint64, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	todoID := uint(idUint64)
	todo.ID = uint(todoID)

	err = json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !todo.Priority.IsValid() {
		http.Error(w, "Invalid priority value", http.StatusBadRequest)
		return
	}
	if !todo.Category.IsValid() {
		http.Error(w, "Invalid category value", http.StatusBadRequest)
		return
	}

	err = todo.UpdateTodo(todoID, userID)
	if err != nil {
		http.Error(w, "Error trying to update todo", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Todo updated successfully!",
		"todo":    todo,
	})
}

// DeleteTodo godoc
// @Summary      Delete a todo by ID
// @Description  Delete a todo task by its ID
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id  path     int  true  "Todo ID"
// @Success      200  {string}  string  "Todo successfully deleted"
// @Failure      400  {string}  string  "Invalid ID"
// @Failure      401  {string}  string  "Unauthorized - missing token"
// @Failure      500  {string}  string  "Error trying to delete todo"
// @Security BearerToken
// @Router       /tasks/{id} [delete]
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	var todo models.Todo
	todo.UserID = userID

	todoIDStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	idUint64, err := strconv.ParseUint(todoIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	todoID := uint(idUint64)
	todo.ID = uint(todoID)

	err = todo.DeleteTodo(todoID, userID)
	if err != nil {
		http.Error(w, "Error trying to delete todo", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Todo deleted successfully!",
	})
}

func getUserIdFromToken(r *http.Request) (uint, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return 0, http.ErrNoCookie
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, http.ErrNoCookie
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, http.ErrNoCookie
	}

	userID := uint(claims["user_id"].(float64))
	return userID, nil
}

func parsePaginationParams(r *http.Request) (page, limit, offset int) {
	page = 1
	limit = 10

	if pStr := r.URL.Query().Get("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset = (page - 1) * limit
	return
}

func parseSortParams(r *http.Request) (sortBy, sortOrder string) {
	sortBy = r.URL.Query().Get("sort_by")
	sortOrder = r.URL.Query().Get("order")

	if sortBy == "" {
		sortBy = "id"
	}

	if sortOrder == "" {
		sortOrder = "asc"
	}

	validSorts := map[string]bool{
		"id":        true,
		"complete":  true,
		"priority":  true,
		"category":  true,
		"createdat": true,
		"duedate":   true,
	}

	if !validSorts[sortBy] {
		sortBy = "id"
	}

	return sortBy, sortOrder
}

func parseFilterParams(r *http.Request) (priority models.Priorities, category models.Categories, complete *bool) {
	priority = models.Priorities(r.URL.Query().Get("priority"))
	category = models.Categories(r.URL.Query().Get("category"))

	if completeStr := r.URL.Query().Get("complete"); completeStr != "" {
		val := completeStr == "true"
		complete = &val
	}

	return
}
