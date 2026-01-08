package repository

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func AddInstance(userId uint, api models.Plugin, url string) error {
	return database.DB.Create(&models.Instance{
		UserId:   userId,
		Api:      api.Name(),
		Provider: api.Provider(),
		Url:      url,
	}).Error
}

func ListInstances() ([]models.Instance, error) {
	var instances []models.Instance
	err := database.DB.Find(&instances).Error
	return instances, err
}

func ListInstancesByUserID(userId uint) ([]models.Instance, error) {
	var instances []models.Instance
	err := database.DB.Find(&instances, "user_id = ?", userId).Error
	return instances, err
}

func ListInstancesByUserIDByAPI(userId uint, api string) ([]models.Instance, error) {
	var instances []models.Instance
	err := database.DB.Find(&instances, "user_id = ? AND api = ?", userId, api).Error
	return instances, err
}

func ListInstancesByUserIDByProvider(userId uint, provider string) ([]models.Instance, error) {
	var instances []models.Instance
	err := database.DB.Find(&instances, "user_id = ? AND provider = ?", userId, provider).Error
	return instances, err
}

func DeleteInstance(id uint) error {
	return database.DB.Delete(&models.Instance{}, id).Error
}

func DeleteInstanceByUserID(userId uint, id uint) error {
	return database.DB.Delete(&models.Instance{}, "user_id = ? AND id = ?", userId, id).Error
}
