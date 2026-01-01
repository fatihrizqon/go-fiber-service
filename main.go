package main

import (
	"fmt"
	"log"

	_ "github.com/fatihrizqon/go-fiber-service/docs"
	"github.com/fatihrizqon/go-fiber-service/logger"
	"github.com/fatihrizqon/go-fiber-service/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @title Go REST API with Fiber Framework
// @version 1.0
// @description This is an Official Documentation for Go REST API with Fiber Framework
// @termsOfService http://swagger.io/terms/
// @contact.name Fatih Rizqon
// @contact.email fatihrizqon@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:3000
// @BasePath /
func main() {
	logger.Init()

	if logLevel := "info"; logLevel != "" {
		if err := logger.SetLogLevel(logLevel); err != nil {
			log.Fatalf("Invalid log level: %s", logLevel)
		}
	}

	fmt.Println("Starting the server...")
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	router.NewRouter(app)
	log.Fatal(app.Listen("127.0.0.1:3000"))
}
