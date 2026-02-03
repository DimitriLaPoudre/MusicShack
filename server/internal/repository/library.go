package repository

import (
	"fmt"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func ListSong(limit uint, offset uint) ([]models.Song, error) {
	var songs []models.Song
	if err := database.DB.
		Limit(int(limit)).
		Offset(int(offset)).
		Order("path ASC").
		Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("repository.ListSong: %w", err)
	}
	return songs, nil
}

func ListSongByUserID(userId uint, limit uint, offset uint) ([]models.Song, error) {
	var songs []models.Song

	if err := database.DB.
		Where("user_id = ?", userId).
		Limit(int(limit)).
		Offset(int(offset)).
		Order("path ASC").
		Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("repository.ListSongByUserID: %w", err)
	}

	return songs, nil
}

func SaveSong(song models.Song) error {
	if err := database.DB.Save(&song).Error; err != nil {
		return fmt.Errorf("repository.SaveSong: %w", err)
	}
	return nil
}

func SaveSongByUserID(userId uint, song models.Song) error {
	if err := database.DB.
		Model(&models.Song{}).
		Where("id = ? AND user_id = ?", song.ID, userId).
		Updates(song).Error; err != nil {
		return fmt.Errorf("repository.SaveSongByUserID: %w", err)
	}
	return nil
}

func DeleteSong(id uint) error {
	if err := database.DB.Delete(&models.Song{}, id).Error; err != nil {
		return fmt.Errorf("repository.DeleteSong: %w", err)
	}
	return nil
}

func DeleteSongByUserID(userId uint, id uint) error {
	if err := database.DB.
		Where("id = ? AND user_id = ?", id, userId).
		Delete(&models.Song{}).Error; err != nil {
		return fmt.Errorf("repository.DeleteSongByUserID: %w", err)
	}
	return nil
}
