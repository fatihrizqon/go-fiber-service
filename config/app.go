package config

import (
	"github.com/fatihrizqon/go-fiber-service/internal/delivery/handler"
	"github.com/fatihrizqon/go-fiber-service/internal/delivery/http/route"
	"github.com/fatihrizqon/go-fiber-service/internal/repository"
	"github.com/fatihrizqon/go-fiber-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	JWT      *JWTService
}

func Bootstrap(config *BootstrapConfig) {
	// Initialize repositories
	userRepository := repository.NewUserRepository(config.DB)
	authRepository := repository.NewAuthRepository(config.DB)

	// Initialize services
	userService := service.NewUserService(userRepository, config.Validate)
	authService := service.NewAuthService(authRepository, config.Validate)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	routeConfig := route.RouteConfig{
		App:         config.App,
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}

	routeConfig.Setup()
}
