package dto

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type Follow struct {
	ID               uuid.UUID `db:"id"`
	UserID           uuid.UUID `db:"user_id"`
	Provider         string    `db:"provider"`
	ArtistID         string    `db:"artist_id"`
	ArtistName       string    `db:"artist_name"`
	ArtistPictureURL string    `db:"artist_picture_url"`
	Featuring        bool      `db:"featuring"`
}

func (f Follow) ToFollow() model.Follow {
	return model.Follow{
		ID:               f.ID,
		UserID:           f.UserID,
		Provider:         f.Provider,
		ArtistID:         f.ArtistID,
		ArtistName:       f.ArtistName,
		ArtistPictureURL: f.ArtistPictureURL,
		Featuring:        f.Featuring,
	}
}

func FollowsToFollows(dto []Follow) []model.Follow {
	follows := []model.Follow{}
	for _, f := range dto {
		follows = append(follows, f.ToFollow())
	}
	return follows
}
