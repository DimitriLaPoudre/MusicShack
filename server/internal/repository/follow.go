package repository

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func AddFollow(userId uint, provider string, artistId string, artistName string, artistPictureUrl string) error {
	return database.DB.Create(&models.Follow{
		UserId:           userId,
		Provider:         provider,
		ArtistId:         artistId,
		ArtistName:       artistName,
		ArtistPictureUrl: artistPictureUrl,
	}).Error
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

func DeleteFollow(id uint) error {
	return database.DB.Delete(&models.Follow{}, id).Error
}

func DeleteFollowByUserID(userId uint, id uint) error {
	return database.DB.Delete(&models.Follow{}, "user_id = ? AND id = ?", userId, id).Error
}
