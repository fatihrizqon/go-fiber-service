package util

import (
	"sync"
	"time"

	"github.com/fatihrizqon/go-fiber-service/internal/entity"
	"github.com/golang-jwt/jwt/v4"
)

type JWTService struct {
	secret        string
	refreshSecret string
	expiration    time.Duration
}

func (j *JWTService) CreateToken(user *entity.User) (string, error) {
	secret := []byte(j.secret)

	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (j *JWTService) RefreshToken(user *entity.User) (string, error) {
	refresh_secret := []byte(j.refreshSecret)

	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refresh_secret)
}

func (j *JWTService) ParseToken(token string) (jwt.MapClaims, error) {
	jwtSecret := []byte(j.secret)

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
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
