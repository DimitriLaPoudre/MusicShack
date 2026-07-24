package request

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type CreateFollow struct {
	Provider  string `json:"provider" binding:"required"`
	ArtistID  string `json:"artist_id" binding:"required"`
	Featuring bool   `json:"featuring"`
}

func (req CreateFollow) IntoFollow(userID uuid.UUID) (model.Follow, error) {
	return model.Follow{
		UserID:    userID,
		Provider:  req.Provider,
		ArtistID:  req.ArtistID,
		Featuring: req.Featuring,
	}, nil
}
