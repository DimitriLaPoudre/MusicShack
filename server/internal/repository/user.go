package repository

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func ListUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	return users, err
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.Take(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := database.DB.Take(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id uint, updates *models.RequestUser) error {
	result := database.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func DeleteUser(id uint) error {
	return database.DB.Delete(&models.User{}, id).Error
}
