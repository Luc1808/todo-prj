package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Luc1808/todo-prj/internal/models"
	"github.com/Luc1808/todo-prj/internal/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = user.Save()
	if err != nil {
		http.Error(w, "Problems registering user", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User succesfully created!")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = user.VerifyCredentials()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenPair, err := utils.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		http.Error(w, "Unable to create athentication token", http.StatusInternalServerError)
		return
	}

	err = models.StoreRefreshToken(user.ID, tokenPair.RefreshToken, time.Now().Add(time.Hour*24*7))
	if err != nil {
		http.Error(w, "Unable to store athentication token in DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Login successful",
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
