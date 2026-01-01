package router

import (
	"log"

	"github.com/fatihrizqon/go-fiber-service/config"
	"github.com/fatihrizqon/go-fiber-service/internal/handler"
	"github.com/fatihrizqon/go-fiber-service/internal/repository"
	"github.com/fatihrizqon/go-fiber-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App) {
	env, err := config.DotEnv()

	if err != nil {
		log.Fatalln("could not load environment variables", err)
	}

	// establish database connection
	db := config.ConnectDatabase(&env)
	validate := validator.New()

	// run auto migrate
	Migrate(db)

	// Register the Repositories
	userRepository := repository.NewUserRepository(db)
	authRepository := repository.NewAuthRepository(db)

	// Register the Services
	userService := service.NewUserService(userRepository, validate)
	authService := service.NewAuthService(authRepository, validate)

	// Register the Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  200,
			"message": "Go REST API with Fiber Framework",
		})
	})

	// app.Post("/api/v1/auth/register", authHandler.Register) // unimplemented
	app.Post("/api/v1/auth/login", authHandler.Login)
	app.Post("/api/v1/auth/refresh", authHandler.Refresh)
	app.Post("/api/v1/auth/logout", authHandler.Logout)
	app.Get("/api/v1/auth/me", authHandler.Me)

	/*
	 * Wrapping in JWT Middleware
	 */
	// api := app.Group("/api/v1", middleware.JWT)
	api := app.Group("/api/v1")

	api.Post("/users", userHandler.Create)
	api.Get("/users", userHandler.FindAll)
	api.Get("/users/:id", userHandler.FindById)
	api.Put("/users/:id", userHandler.Update)
	api.Delete("/users/:id", userHandler.Delete)

	// api.Post("/logout", func(c *fiber.Ctx) error {
	// 	token := c.Get("Authorization")
	// 	helper.BlacklistToken(token)
	// 	return c.Status(200).JSON(fiber.Map{
	// 		"status":  200,
	// 		"message": "you are unauthenticated",
	// 	})
	// })
}
