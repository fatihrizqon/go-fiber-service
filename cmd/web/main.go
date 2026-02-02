package main

import (
	"fmt"

	"github.com/fatihrizqon/go-fiber-service/config"
	"github.com/fatihrizqon/go-fiber-service/database"
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
	viper := config.NewViper()
	log := config.NewLogger(viper)
	db := config.NewDatabase(viper, log)
	validate := config.NewValidator(viper)
	app := config.NewFiber(viper, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viper,
	})

	// Run database migration
	database.Migrate(db)

	port := viper.GetInt("web.port")

	err := app.Listen(fmt.Sprintf(":%d", port))

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
