package helper

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatihrizqon/go-fiber-service/internal/entity"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var refreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))

func GenerateAccessToken(user entity.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	fmt.Println(claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(user entity.User) (string, error) {
	fmt.Println(user)
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

func ParseToken(tokenString string, isRefresh bool) (jwt.MapClaims, error) {
	secret := jwtSecret
	if isRefresh {
		secret = refreshSecret
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || time.Now().Unix() > int64(exp) {
		return nil, fmt.Errorf("token has expired")
	}

	if _, ok := claims["username"]; !ok {
		claims["username"] = ""
	}
	if _, ok := claims["name"]; !ok {
		claims["name"] = ""
	}

	return claims, nil
}

var TokenBlacklist = struct {
	sync.RWMutex
	tokens map[string]struct{}
}{tokens: make(map[string]struct{})}

func BlacklistToken(token string) {
	TokenBlacklist.Lock()
	defer TokenBlacklist.Unlock()
	TokenBlacklist.tokens[token] = struct{}{}
}

func IsBlacklisted(token string) bool {
	TokenBlacklist.RLock()
	defer TokenBlacklist.RUnlock()
	_, exists := TokenBlacklist.tokens[token]
	return exists
}
