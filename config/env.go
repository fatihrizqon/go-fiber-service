package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	host       string `mapstructure:"DATABASE_HOST"`
	port       string `mapstructure:"DATABASE_PORT"`
	database   string `mapstructure:"DATABASE_NAME"`
	username   string `mapstructure:"DATABASE_USER"`
	password   string `mapstructure:"DATABASE_PASSWORD"`
	jwt_secret string `mapstructure:"JWT_SECRET"`
}

func DotEnv() (env Environment, err error) {
	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	env.host = os.Getenv("DATABASE_HOST")
	env.port = os.Getenv("DATABASE_PORT")
	env.database = os.Getenv("DATABASE_NAME")
	env.username = os.Getenv("DATABASE_USER")
	env.password = os.Getenv("DATABASE_PASSWORD")
	env.jwt_secret = os.Getenv("JWT_SECRET")

	return
}
