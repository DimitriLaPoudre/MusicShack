package model

import (
	"context"
	"io"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByFilter(ctx context.Context, filter UserFilter) (User, error)
	ListUsersByFilter(ctx context.Context, filter UserFilter) ([]User, error)
	UpdateUser(ctx context.Context, partialUser PartialUser) (User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type SessionRepository interface {
	CreateSession(ctx context.Context, s Session) (Session, error)
	GetSessionByFilter(ctx context.Context, filter SessionFilter) (Session, error)
	ListSessionsByFilter(ctx context.Context, filter SessionFilter) ([]Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteSessionByToken(ctx context.Context, token string) error
	DeleteSessionExpired(ctx context.Context) error
}

type InstanceRepository interface {
	CreateInstance(ctx context.Context, i Instance) (Instance, error)
	ListInstancesByFilter(ctx context.Context, filter InstanceFilter) ([]Instance, error)
	DeleteInstance(ctx context.Context, id uuid.UUID) error
	DeleteInstanceByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type Migrator interface {
	Migrate(dsn string) error
}

type Plugin interface {
	Name() string
	Provider() string
	Status(ctx context.Context, url string) error
	Download(ctx context.Context, user *User, id string) (io.ReadCloser, string, error)
	Song(ctx context.Context, instances []Instance, id string) (Song, error)
	Album(ctx context.Context, instances []Instance, id string) (Album, error)
	Artist(ctx context.Context, instances []Instance, id string) (Artist, error)
	Playlist(ctx context.Context, instances []Instance, id string) (Playlist, error)
	Search(ctx context.Context, instances []Instance, song string, album string, artist string, playlist string) (Search, error)
	Url(ctx context.Context, instances []Instance, url string) (UrlItem, error)
	Lyrics(ctx context.Context, instances []Instance, id string) (string, string, error)
}
