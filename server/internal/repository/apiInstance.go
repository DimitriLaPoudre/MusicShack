package repository

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func AddInstance(userId uint, api string, url string) error {
	return database.DB.Create(&models.ApiInstance{
		UserId: userId,
		Api:    api,
		Url:    url,
	}).Error
}

func ListInstances() ([]models.ApiInstance, error) {
	var instances []models.ApiInstance
	err := database.DB.Find(&instances).Error
	return instances, err
}

func ListInstancesByUserID(userId uint) ([]models.ApiInstance, error) {
	var instances []models.ApiInstance
	err := database.DB.Find(&instances, "user_id = ?", userId).Error
	return instances, err
}

func ListInstancesByUserIDByAPI(userId uint, api string) ([]models.ApiInstance, error) {
	var instances []models.ApiInstance
	err := database.DB.Find(&instances, "user_id = ? AND api = ?", userId, api).Error
	return instances, err
}

func DeleteInstance(id uint) error {
	return database.DB.Delete(&models.ApiInstance{}, id).Error
}

func DeleteInstanceByUserID(userId uint, id uint) error {
	return database.DB.Delete(&models.ApiInstance{}, "user_id = ? AND id = ?", userId, id).Error
}
