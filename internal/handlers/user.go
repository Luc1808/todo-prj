package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Luc1808/todo-prj/internal/models"
	"github.com/Luc1808/todo-prj/internal/utils"
	"github.com/golang-jwt/jwt/v5"
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

	err = models.StoreRefreshToken(user.ID, tokenPair.RefreshToken, time.Now().Add(utils.RefreshTokenDuration))
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

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var tokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedToken, err := jwt.Parse(tokenRequest.RefreshToken, func(token *jwt.Token) (any, error) {
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
	userEmail := claims["email"].(string)

	storedToken, err := models.FindRefreshToken(tokenRequest.RefreshToken, userID)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}
	newTokenPair, err := utils.GenerateTokenPair(userID, userEmail)
	if err != nil {
		http.Error(w, "Failed to generate new token pair", http.StatusInternalServerError)
		return
	}

	err = models.StoreRefreshToken(userID, newTokenPair.RefreshToken, time.Now().Add(utils.RefreshTokenDuration))
	if err != nil {
		http.Error(w, "Failed to save new refresh token", http.StatusInternalServerError)
		return
	}

	err = models.DeleteRefreshToken(storedToken.Token)
	if err != nil {
		http.Error(w, "Failed to delete old refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":       "Refresh successful. New token pair added!",
		"access_token":  newTokenPair.AccessToken,
		"refresh_token": newTokenPair.RefreshToken,
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
