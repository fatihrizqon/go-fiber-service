package middleware

import (
	"strings"

	"github.com/fatihrizqon/go-fiber-service/internal/delivery/http/request"
	"github.com/fatihrizqon/go-fiber-service/internal/util"
	"github.com/gofiber/fiber/v2"
)

func NewAuth() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization", "")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "authorization header missing",
				"status":  fiber.StatusUnauthorized,
			})
		}

		// Expect: "Bearer <token>"
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid authorization format",
				"status":  fiber.StatusUnauthorized,
			})
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		req := &request.VerifyUserRequest{
			Token: token,
		}

		auth, err := util.ParseToken(req.Token)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized: " + err.Error(),
				"status":  fiber.StatusUnauthorized,
			})
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}
