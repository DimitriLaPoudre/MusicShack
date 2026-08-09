package dto

import (
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type Song struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	ISRC      string    `db:"isrc"`
	Path      string    `db:"path"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s Song) ToSong() model.Song {
	return model.Song{
		ID:        s.ID,
		UserID:    s.UserID,
		ISRC:      s.ISRC,
		Path:      s.Path,
		UpdatedAt: s.UpdatedAt,
	}
}

func SongsToSongs(dto []Song) []model.Song {
	songs := []model.Song{}
	for _, s := range dto {
		songs = append(songs, s.ToSong())
	}
	return songs
}
