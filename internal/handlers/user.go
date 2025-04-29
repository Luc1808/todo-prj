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

// RegisterHandler godoc
// @Summary      Register a new user
// @Description  Register a new user with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body     models.UserRequest  true  "User object"
// @Success      201  {string}  string  "User successfully created"
// @Failure      400  {string}  string  "Invalid JSON"
// @Failure      500  {string}  string  "Problems registering user"
// @Router       /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = user.Save()
	if err != nil {
		http.Error(w, "Problems registering user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User succesfully created!")
}

// LoginHandler godoc
// @Summary      Login user
// @Description  Login user with email and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body     models.UserRequest  true  "User object"
// @Success      200  {string}  string  "Login successful"
// @Failure      400  {string}  string  "Invalid JSON"
// @Failure      401  {string}  string  "Invalid credentials"
// @Failure      500  {string}  string  "Unable to create authentication token"
// @Router       /login [post]
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

// RefreshHandler godoc
// @Summary      Refresh access token
// @Description  Refresh access token using refresh token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        refresh_token  body     models.RefreshTokenRequest true  "Refresh token"
// @Success      200  {string}  string  "New access token"
// @Failure      400  {string}  string  "Invalid JSON"
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Failed to generate new token pair"
// @Router       /refresh [post]
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

// func GetUsers(w http.ResponseWriter, r *http.Request) {
// 	users, err := models.GetAllUsers()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(users)
// }
