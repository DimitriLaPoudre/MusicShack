package model

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByFilter(ctx context.Context, filter FilterUser) (User, error)
	ListUsersByFilter(ctx context.Context, filter FilterUser) ([]User, error)
	UpdateUser(ctx context.Context, partialUser PartialUser) (User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, s Session) (Session, error)
	GetSessionByID(ctx context.Context, sessionID uuid.UUID) (Session, error)
	GetSessionByToken(ctx context.Context, token string) (Session, error)
	GetSessionByUserID(ctx context.Context, userID uuid.UUID) ([]Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteSessionByToken(ctx context.Context, token string) error
}

type InstanceRepository interface {
	CreateInstance(ctx context.Context, i Instance) (Instance, error)
	GetInstance(ctx context.Context, id uuid.UUID) (Instance, error)
	ListInstancesByFilter(ctx context.Context, filter InstanceFilter) ([]Instance, error)
	DeleteInstance(ctx context.Context, id uuid.UUID) error
	DeleteInstanceByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type Plugin interface {
	Name() string
	Provider() string
	Status(ctx context.Context, url string) error
	Download(ctx context.Context, user *User, id string) (io.ReadCloser, string, error)
	Song(ctx context.Context, url string, id string) (Song, error)
	Playlist(ctx context.Context, url string, id string) (Playlist, error)
	Album(ctx context.Context, url string, id string) (Album, error)
	Artist(ctx context.Context, url string, id string) (Artist, error)
	Search(ctx context.Context, url string, song string, album string, artist string) (Search, error)
	Url(ctx context.Context, url string, id string) (UrlItem, error)
	Lyrics(ctx context.Context, url string, id string) (string, string, error)
}
