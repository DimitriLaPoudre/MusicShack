package database

import (
	"fmt"
	"log"
	"os"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

func initAdmin() (*models.Admin, error) {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		log.Fatal("ADMIN_PASSWORD is missing")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("initAdmin: %w", err)
	}

	admin := models.Admin{
		Password: string(hashedPassword),
	}
	return &admin, nil
}

func migrationDB(db *gorm.DB) error {
	m := db.Migrator()

	if m.HasIndex(&models.Song{}, "idx_songs_user_id") {
		if err := m.DropIndex(&models.Song{}, "idx_songs_user_id"); err != nil {
			return fmt.Errorf("drop idx_songs_user_id: %w", err)
		}
	}

	if m.HasIndex(&models.Song{}, "idx_songs_isrc") {
		if err := m.DropIndex(&models.Song{}, "idx_songs_isrc"); err != nil {
			return fmt.Errorf("drop idx_songs_isrc: %w", err)
		}
	}

	if !m.HasIndex(&models.Song{}, "idx_song") {
		if err := m.CreateIndex(&models.Song{}, "idx_song"); err != nil {
			return fmt.Errorf("create idx_song: %w", err)
		}
	}

	return nil
}

func init() {
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		postgresHost, postgresUser, postgresPassword, postgresDB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.UserSession{}, &models.Instance{}, &models.Follow{}, &models.Song{}, &models.Admin{}); err != nil {
		log.Fatal(err)
	}

	if err := migrationDB(db); err != nil {
		log.Fatal(err)
	}

	admin, err := initAdmin()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(admin).Error
	if err != nil {
		log.Fatal(err)
	}

	DB = db
}
