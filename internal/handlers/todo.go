package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Luc1808/todo-prj/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID for todo creation
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

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

func GetAllTodosWithPagination(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

	// Pagination, sorting stuff
	totalTodos, err := models.GetNumberOfTodos(userID)
	if err != nil {
		http.Error(w, "Failed to get the total number of to-dos", http.StatusUnauthorized)
		return
	}

	page, limit, offset := parsePaginationParams(r)
	sortBy, sortOrder := parseSortParams(r)

	todos, err := models.GetTodosWithPagination(userID, limit, offset, sortBy, sortOrder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get total pages
	totalPages := int(math.Ceil(float64(totalTodos) / float64(limit)))

	response := map[string]any{
		"page":        page,
		"limit":       limit,
		"total_todos": totalTodos,
		"total_pages": totalPages,
		"sort_by":     sortBy,
		"order":       sortOrder,
		"todos":       todos,
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetTodoByID(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

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

func GetTodosByPriority(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

	todoPriority := models.Priorities(strings.TrimPrefix(r.URL.Path, "/tasks/priority/"))

	todos, err := models.GetTodosByPriority(userID, todoPriority)
	if err != nil {
		http.Error(w, "Error trying to get all todos", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func GetTodosByCategory(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

	todoCategory := models.Categories(strings.TrimPrefix(r.URL.Path, "/tasks/category/"))

	todos, err := models.GetTodosByCategory(userID, todoCategory)
	if err != nil {
		http.Error(w, "Error trying to get all todos", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func GetCompletedTodos(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

	todos, err := models.GetCompletedTodos(userID)
	if err != nil {
		http.Error(w, "Error trying to get all todos", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

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

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Unauthorized - missing token", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_REFRESH_TOKEN")), nil
	})
	if err != nil || !parsedToken.Valid {
		http.Error(w, "Failed to parse token", http.StatusUnauthorized)
		return
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Failed to parse token claims", http.StatusUnauthorized)
		return
	}

	userID := uint(claims["user_id"].(float64))

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
