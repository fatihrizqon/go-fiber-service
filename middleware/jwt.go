package middleware

import (
	"net/http"
	"os"

	"github.com/fatihrizqon/go-fiber-service/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func JWT(c *fiber.Ctx) error {
	token := c.Cookies("access_token")

	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "no token provided"})
	}

	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(http.StatusUnauthorized, "invalid signing method")
		}
		return jwtSecret, nil
	})

	if err != nil || !parsedToken.Valid || helper.IsBlacklisted(token) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}

	if claims.Id == "" || claims.Username == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
	}

	c.Locals("id", claims.Id)

	return c.Next()
}
