package repository

import (
	"fmt"
	"time"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func GetAdmin() (*models.Admin, error) {
	var admin models.Admin
	err := database.DB.First(&admin, "id = true").Error
	if err != nil {
		return nil, fmt.Errorf("repository.GetAdmin: %w", err)
	}
	return &admin, nil
}

func ChangeAdminPassword(password string) error {
	err := database.DB.Model(&models.Admin{}).Updates(map[string]any{"password": password, "token": ""}).Error
	if err != nil {
		return fmt.Errorf("repository.ChangeAdminPassword: %w", err)
	}
	return nil
}

func CreateAdminSession(token string, expiresAt time.Time) error {
	err := database.DB.Model(&models.Admin{}).Updates(map[string]any{"token": token, "expires_at": expiresAt}).Error
	if err != nil {
		return fmt.Errorf("repository.CreateAdminSession: %w", err)
	}
	return nil
}

func DeleteAdminSession() error {
	err := database.DB.Model(&models.Admin{}).Update("token", "").Error
	if err != nil {
		return fmt.Errorf("repository.DeleteAdminSession: %w", err)
	}
	return nil
}
