package repository

import (
	"fmt"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func AddInstance(userId uint, api models.Plugin, url string) error {
	if err := database.DB.Create(&models.Instance{
		UserId:   userId,
		Api:      api.Name(),
		Provider: api.Provider(),
		Url:      url,
	}).Error; err != nil {
		return fmt.Errorf("repository.AddInstance: %w", err)
	}
	return nil
}

func ListInstances() ([]models.Instance, error) {
	var instances []models.Instance
	if err := database.DB.Find(&instances).Error; err != nil {
		return nil, fmt.Errorf("repository.ListInstances: %w", err)
	}
	return instances, nil
}

func ListInstancesByUserID(userId uint) ([]models.Instance, error) {
	var instances []models.Instance
	if err := database.DB.Find(&instances, "user_id = ?", userId).Error; err != nil {
		return nil, fmt.Errorf("repository.ListInstancesByUserID: %w", err)
	}
	return instances, nil
}

func ListInstancesByUserIDByAPI(userId uint, api string) ([]models.Instance, error) {
	var instances []models.Instance
	if err := database.DB.Find(&instances, "user_id = ? AND api = ?", userId, api).Error; err != nil {
		return nil, fmt.Errorf("repository.ListInstancesByUserIDByAPI: %w", err)
	}
	return instances, nil
}

func ListInstancesByUserIDByProvider(userId uint, provider string) ([]models.Instance, error) {
	var instances []models.Instance
	if err := database.DB.Find(&instances, "user_id = ? AND provider = ?", userId, provider).Error; err != nil {
		return nil, fmt.Errorf("repository.ListInstancesByUserIDByProvider: %w", err)
	}
	return instances, nil
}

func DeleteInstance(id uint) error {
	if err := database.DB.Delete(&models.Instance{}, id).Error; err != nil {
		return fmt.Errorf("repository.DeleteInstance: %w", err)
	}
	return nil
}

func DeleteInstanceByUserID(userId uint, id uint) error {
	if err := database.DB.Delete(&models.Instance{}, "user_id = ? AND id = ?", userId, id).Error; err != nil {
		return fmt.Errorf("repository.DeleteInstanceByUserID: %w", err)
	}
	return nil
}
