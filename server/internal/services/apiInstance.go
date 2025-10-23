package services

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func ConstructApiInstance(apiInstance *models.ApiInstanceRequest) *models.ApiInstance {
	return &models.ApiInstance{
		Api: apiInstance.Api,
		Url: apiInstance.Url,
	}
}
