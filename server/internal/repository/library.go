package repository

import "github.com/DimitriLaPoudre/MusicShack/server/internal/models"

func AddSong(song models.Song) error {
	return nil
}

func AddSongByUserID(userId uint, song models.Song) error {
	return nil
}

func ListSong(limit uint, offset uint) ([]models.Song, error) {
	return []models.Song{}, nil
}

func ListSongByUserID(userId uint, limit uint, offset uint) ([]models.Song, error) {
	return []models.Song{}, nil
}

func UpdateSong(song models.Song) error {
	return nil
}

func UpdateSongByUserID(userId uint, song models.Song) error {
	return nil
}

func DeleteSong(id uint) error {
	return nil
}

func DeleteSongByUserID(userId uint, id uint) error {
	return nil
}
