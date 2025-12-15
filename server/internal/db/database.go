package database

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("./db/data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.ApiInstance{}, &models.Follow{})
	DB = db
}
