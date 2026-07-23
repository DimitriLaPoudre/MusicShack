package model

import "github.com/google/uuid"

type Follow struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	Provider      string
	ArtistID      string
	ArtistName    string
	ArtistPicture string
	Featuring     bool
}

type FollowFilter struct {
	ID            *uuid.UUID
	UserID        *uuid.UUID
	Provider      *string
	ArtistID      *string
	ArtistName    *string
	ArtistPicture *string
	Featuring     *bool
}
