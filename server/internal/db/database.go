package database

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {
	db, err := gorm.Open(sqlite.Open("./db/data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.ApiInstance{})
	DB = db
}
