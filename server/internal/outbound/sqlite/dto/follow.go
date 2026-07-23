package dto

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type Follow struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Provider      string    `json:"provider"`
	ArtistID      string    `json:"artist_id"`
	ArtistName    string    `json:"artist_name"`
	ArtistPicture string    `json:"artist_picture"`
	Featuring     bool      `json:"featuring"`
}

func (f Follow) ToFollow() model.Follow {
	return model.Follow{
		ID:            f.ID,
		UserID:        f.UserID,
		Provider:      f.Provider,
		ArtistID:      f.ArtistID,
		ArtistName:    f.ArtistName,
		ArtistPicture: f.ArtistPicture,
		Featuring:     f.Featuring,
	}
}

func FollowsToFollows(dto []Follow) []model.Follow {
	follows := []model.Follow{}
	for _, f := range dto {
		follows = append(follows, f.ToFollow())
	}
	return follows
}
