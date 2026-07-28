package model

import (
	"context"
	"io"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	TransactionRepository
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserByFilter(ctx context.Context, filter UserFilter) (User, error)
	ListUsersByFilter(ctx context.Context, filter UserFilter) ([]User, error)
	UpdateUser(ctx context.Context, partialUser PartialUser) (User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type SessionRepository interface {
	TransactionRepository
	CreateSession(ctx context.Context, s Session) (Session, error)
	GetSessionByFilter(ctx context.Context, filter SessionFilter) (Session, error)
	ListSessionsByFilter(ctx context.Context, filter SessionFilter) ([]Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteSessionByToken(ctx context.Context, token string) error
	DeleteSessionExpired(ctx context.Context) error
}

type InstanceRepository interface {
	TransactionRepository
	CreateInstance(ctx context.Context, i Instance) (Instance, error)
	ListInstancesByFilter(ctx context.Context, filter InstanceFilter) ([]Instance, error)
	UpdateInstancePing(ctx context.Context, id uuid.UUID, ping *time.Duration) (Instance, error)
	DeleteInstance(ctx context.Context, id uuid.UUID) error
	DeleteInstanceByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type FollowRepository interface {
	TransactionRepository
	CreateFollow(ctx context.Context, f Follow) (Follow, error)
	GetFollowByFilter(ctx context.Context, filter FollowFilter) (Follow, error)
	ListFollowsByFilter(ctx context.Context, filter FollowFilter) ([]Follow, error)
	DeleteFollow(ctx context.Context, id uuid.UUID) error
	DeleteFollowByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type TransactionRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Migrator interface {
	Migrate(dsn string) error
}

type Plugin interface {
	Name() string
	Provider() string
	Status(ctx context.Context, url string) error
	Download(ctx context.Context, instances []Instance, id string, hiRes bool) (io.ReadCloser, string, error)
	SongInfo(ctx context.Context, instances []Instance, id string) (SongInfo, error)
	SongInfoByISRC(ctx context.Context, instances []Instance, isrc string) (SongInfo, error)
	AlbumInfo(ctx context.Context, instances []Instance, id string) (AlbumInfo, error)
	AlbumSongs(ctx context.Context, instances []Instance, id string, limit int, offset int) (PaginatedSongs, error)
	ArtistInfo(ctx context.Context, instances []Instance, id string) (ArtistInfo, error)
	ArtistAlbums(ctx context.Context, instances []Instance, id string, limit int, offset int) (ArtistPaginatedAlbums, error)
	PlaylistInfo(ctx context.Context, instances []Instance, id string) (PlaylistInfo, error)
	PlaylistSongs(ctx context.Context, instances []Instance, id string, limit int, offset int) (PaginatedSongs, error)
	Search(ctx context.Context, instances []Instance, song string, album string, artist string, playlist string, limit int, offset int) (Search, error)
	URL(ctx context.Context, url string) (URLItem, error)
}
