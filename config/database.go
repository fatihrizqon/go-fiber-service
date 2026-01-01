package config

import (
	"fmt"

	"github.com/fatihrizqon/go-fiber-service/helper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(config *Environment) *gorm.DB {
	credentials := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.host, config.port, config.username, config.password, config.database)

	db, err := gorm.Open(postgres.Open(credentials), &gorm.Config{})
	helper.PanicIfError(err)

	fmt.Println("Database connection has been established.")

	return db
}
