package repository

import (
	"fmt"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func AddFollow(userId uint, provider string, artistId string, artistName string, artistPictureUrl string) (*models.Follow, error) {
	follow := models.Follow{
		UserId:           userId,
		Provider:         provider,
		ArtistId:         artistId,
		ArtistName:       artistName,
		ArtistPictureUrl: artistPictureUrl,
	}
	if err := database.DB.Create(&follow).Error; err != nil {
		return nil, fmt.Errorf("repository.AddFollow: %w", err)
	}
	return &follow, nil
}

func ListFollows() ([]models.Follow, error) {
	var follows []models.Follow
	err := database.DB.Find(&follows).Error
	return follows, err
}

func ListFollowsByUserID(userId uint) ([]models.Follow, error) {
	var follows []models.Follow
	err := database.DB.Find(&follows, "user_id = ?", userId).Error
	return follows, err
}

func GetFollowByProviderByArtistID(provider string, artistId string) (*models.Follow, error) {
	var follow models.Follow
	if err := database.DB.First(&follow, "provider = ? AND artist_id = ?", provider, artistId).Error; err != nil {
		return nil, fmt.Errorf("repository.GetFollowByProviderByArtistID: %w", err)
	}
	return &follow, nil
}

func DeleteFollow(id uint) error {
	return database.DB.Delete(&models.Follow{}, id).Error
}

func DeleteFollowByUserID(userId uint, id uint) error {
	return database.DB.Delete(&models.Follow{}, "user_id = ? AND id = ?", userId, id).Error
}
