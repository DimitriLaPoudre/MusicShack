package repository

import (
	"fmt"

	database "github.com/DimitriLaPoudre/MusicShack/server/internal/db"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func GetSong(userId uint, id uint) (models.Song, error) {
	var song models.Song

	if err := database.DB.First(&song, "id = ?", id).Error; err != nil {
		return models.Song{}, fmt.Errorf("repository.GetSong: %w", err)
	}

	return song, nil
}

func GetSongByUserID(userId uint, id uint) (models.Song, error) {
	var song models.Song

	if err := database.DB.
		First(&song, "id = ? AND user_id = ?", id, userId).Error; err != nil {
		return models.Song{}, fmt.Errorf("repository.GetSongByUserID: %w", err)
	}

	return song, nil
}

func CountSong(q string) (int64, error) {
	var total int64
	if err := database.DB.Model(&models.Song{}).
		Where("path ILIKE ?", "%"+q+"%").
		Count(&total).Error; err != nil {
		return 0, fmt.Errorf("repository.CountSong: %w", err)
	}
	return total, nil
}

func CountSongByUserID(userId uint, q string) (int64, error) {
	var total int64
	if err := database.DB.Model(&models.Song{}).
		Where("path ILIKE ? AND user_id = ?", "%"+q+"%", userId).
		Count(&total).Error; err != nil {
		return 0, fmt.Errorf("repository.CountSong: %w", err)
	}
	return total, nil
}

func ListSong(q string, limit int, offset int) ([]models.Song, error) {
	var songs []models.Song

	if err := database.DB.
		Where("path ILIKE ?", "%"+q+"%").
		Limit(limit).
		Offset(offset).
		Order("updated_at DESC NULLS LAST").
		Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("repository.ListSong: %w", err)
	}

	return songs, nil
}

func ListSongByUserID(userId uint, q string, limit int, offset int) ([]models.Song, error) {
	var songs []models.Song

	if err := database.DB.
		Where("path ILIKE ? AND user_id = ?", "%"+q+"%", userId).
		Limit(limit).
		Offset(offset).
		Order("updated_at DESC NULLS LAST").
		Find(&songs).Error; err != nil {
		return nil, fmt.Errorf("repository.ListSongByUserID: %w", err)
	}

	return songs, nil
}

func AddSong(song models.Song) error {
	if err := database.DB.Create(&song).Error; err != nil {
		return fmt.Errorf("repository.AddSong: %w", err)
	}
	return nil
}

func UpdateSong(song models.Song) error {
	if err := database.DB.Model(&models.Song{}).
		Where("id = ?", song.ID).
		Updates(song).Error; err != nil {
		return fmt.Errorf("repository.UpdateSong: %w", err)
	}
	return nil
}

func UpdateSongByUserID(userId uint, song models.Song) error {
	if err := database.DB.Model(&models.Song{}).
		Where("id = ? AND user_id = ?", song.ID, userId).
		Updates(song).Error; err != nil {
		return fmt.Errorf("repository.UpdateSongByUserID: %w", err)
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
