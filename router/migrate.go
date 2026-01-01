package router

import (
	"github.com/fatihrizqon/go-fiber-service/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.User{})
}
