package response

import "github.com/DimitriLaPoudre/MusicShack/internal/model"

type Follow struct {
	ID               string `json:"id"`
	UserID           string `json:"user_id"`
	Provider         string `json:"provider"`
	ArtistID         string `json:"artist_id"`
	ArtistName       string `json:"artist_name"`
	ArtistPictureURL string `json:"artist_picture_url"`
	Featuring        bool   `json:"featuring"`
}

func FollowToResponse(follow model.Follow) Follow {
	return Follow{
		ID:               follow.ID.String(),
		UserID:           follow.UserID.String(),
		Provider:         follow.Provider,
		ArtistID:         follow.ArtistID,
		ArtistName:       follow.ArtistName,
		ArtistPictureURL: follow.ArtistPictureURL,
		Featuring:        follow.Featuring,
	}
}

func FollowsToResponse(follows []model.Follow) []Follow {
	r := []Follow{}
	for _, follow := range follows {
		r = append(r, FollowToResponse(follow))
	}
	return r
}
