package util

import (
	"sync"
	"time"

	"github.com/fatihrizqon/go-fiber-service/internal/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var jwtSecret []byte
var refreshSecret []byte

type Claims struct {
	Id string `json:"id"`
	jwt.RegisteredClaims
}

func NewJWT(v *viper.Viper) {
	access := v.GetString("jwt.secret")
	refresh := v.GetString("jwt.refresh_secret")

	if access == "" || refresh == "" {
		panic("JWT secret is empty")
	}

	jwtSecret = []byte(access)
	refreshSecret = []byte(refresh)
}

func CreateToken(user entity.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func RefreshToken(user entity.User) (string, error) {
	claims := jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"name":     user.Name,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(refreshSecret)
}

func ParseToken(tokenString string) (*entity.User, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	userID, err := uuid.Parse(claims.Id)
	if err != nil {
		return nil, fiber.ErrUnauthorized
	}

	return &entity.User{Id: userID}, nil
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
