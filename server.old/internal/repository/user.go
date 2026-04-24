package repository

import (
	"fmt"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	if err := database.DB.Create(user).Error; err != nil {
		return fmt.Errorf("repository.CreateUser: %w", err)
	}
	return nil
}

func ListUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	if err != nil {
		return []models.User{}, fmt.Errorf("repository.ListUsers: %w", err)
	}
	return users, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Take(&user, id).Error
	if err != nil {
		return nil, fmt.Errorf("repository.GetUserByID: %w", err)
	}
	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Take(&user, "username = ?", username).Error
	if err != nil {
		return nil, fmt.Errorf("repository.GetUserByUsername: %w", err)
	}
	return &user, nil
}

func UpdateUser(id uint, updates *models.RequestUser) error {
	result := database.DB.Model(&models.User{}).Where("id = ?", id).Select("hi_res").Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("repository.UpdateUser: %w", result.Error)
	}
	if result.RowsAffected == 0 && (updates.Username != "" || updates.Password != "") {
		return fmt.Errorf("repository.UpdateUser: %w", gorm.ErrRecordNotFound)
	}
	return nil
}

func DeleteUser(id uint) error {
	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		return fmt.Errorf("repository.CreateUser: %w", err)
	}
	return nil
}
