package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var secretAcessKey = os.Getenv("JWT_ACCESS_SECRET")
var secretRefreshKey = os.Getenv("JWT_REFRESH_SECRET")

const (
	// GO BACK TO 15 MINUTES LATER
	AccessTokenDuration  = time.Hour * 3
	RefreshTokenDuration = time.Hour * 24 * 7
)

// func GenerateToken(id uint, email string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"user_id": id,
// 		"email":   email,
// 		"exp":     time.Now().Add(time.Hour * 2).Unix(),
// 	})

// 	return token.SignedString([]byte(secretAcessKey))
// }

func GenerateTokenPair(id uint, email string) (TokenPair, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(AccessTokenDuration).Unix(),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"exp":     time.Now().Add(RefreshTokenDuration).Unix(),
	})

	accessTokenString, err := accessToken.SignedString([]byte(secretAcessKey))
	if err != nil {
		return TokenPair{}, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(secretRefreshKey))
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secretAcessKey), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		fmt.Errorf("invalid token")
	}

	return nil
}
