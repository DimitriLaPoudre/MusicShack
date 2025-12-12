package repository

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func CreateApiInstance(apiInstance *models.ApiInstance) error {
	return database.DB.Create(apiInstance).Error
}

func ListApiInstances() ([]models.ApiInstance, error) {
	var apiInstances []models.ApiInstance
	err := database.DB.Find(&apiInstances).Error
	return apiInstances, err
}

func ListApiInstancesByApi(api string) ([]models.ApiInstance, error) {
	var apiInstances []models.ApiInstance
	err := database.DB.Find(&apiInstances, "api = ?", api).Error
	return apiInstances, err
}

func GetApiInstanceByID(id uint) (*models.ApiInstance, error) {
	var user models.ApiInstance
	err := database.DB.Take(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetApiInstanceByApi(api string) (*models.ApiInstance, error) {
	var user models.ApiInstance
	err := database.DB.Take(&user, "api = ?", api).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetApiInstanceByURL(url string) (*models.ApiInstance, error) {
	var user models.ApiInstance
	err := database.DB.Take(&user, "url = ?", url).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteApiInstance(id uint) error {
	return database.DB.Delete(&models.ApiInstance{}, id).Error
}
