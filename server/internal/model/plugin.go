package model

import (
	"time"

	"github.com/google/uuid"
)

type DataType string

const (
	TypeSong     DataType = "song"
	TypeAlbum    DataType = "album"
	TypeArtist   DataType = "artist"
	TypePlaylist DataType = "playlist"
)

type Quality struct {
	Name  string
	Color string
}

type Pagination struct {
	Limit              int
	Offset             int
	TotalNumberOfItems int
}

// -- Song -- //

type SongInfo struct {
	ID              string
	Title           string
	Duration        int
	ReplayGain      float64
	Peak            float64
	AlbumReplayGain float64
	AlbumPeak       float64
	ReleaseDate     time.Time
	TrackNumber     int
	VolumeNumber    int
	AudioQuality    Quality
	Popularity      int
	Explicit        bool
	ISRC            string
	Artists         []ArtistInfo
	Album           AlbumInfo
}
type EnrichedSongInfo struct {
	Provider string
	SongInfo
}

type PaginatedSongs struct {
	Pagination
	Songs []SongInfo
}
type EnrichedPaginatedSongs struct {
	Provider string
	PaginatedSongs
	Songs []EnrichedSongInfo
}

// -- Album -- //

type AlbumInfo struct {
	ID            string
	Title         string
	Duration      int
	ReleaseDate   time.Time
	NumberTracks  int
	NumberVolumes int
	CoverURL      string
	AudioQuality  Quality
	Popularity    int
	Explicit      bool
	Artists       []ArtistInfo
}
type EnrichedAlbumInfo struct {
	Provider string
	AlbumInfo
}

type PaginatedAlbums struct {
	Pagination
	Albums []AlbumInfo
}
type EnrichedPaginatedAlbums struct {
	Provider string
	PaginatedAlbums
	Albums []EnrichedAlbumInfo
}

// -- Artist -- //

type ArtistInfo struct {
	ID         string
	Name       string
	PictureURL string
	Popularity int
}
type EnrichedArtistInfo struct {
	Provider string
	Followed *uuid.UUID
	ArtistInfo
}

type PaginatedArtists struct {
	Pagination
	Artists []ArtistInfo
}
type EnrichedPaginatedArtists struct {
	Provider string
	PaginatedArtists
	Artists []EnrichedArtistInfo
}

type ArtistPaginatedAlbums struct {
	Albums  PaginatedAlbums
	EPs     PaginatedAlbums
	Singles PaginatedAlbums
}
type EnrichedArtistPaginatedAlbums struct {
	Provider string
	ArtistPaginatedAlbums
	Albums  EnrichedPaginatedAlbums
	EPs     EnrichedPaginatedAlbums
	Singles EnrichedPaginatedAlbums
}

// -- Playlist -- //

type PlaylistInfo struct {
	ID             string
	Title          string
	Description    string
	Duration       int
	LastUpdated    time.Time
	NumberOfTracks int
	CoverURL       string
	Popularity     int
}
type EnrichedPlaylistInfo struct {
	Provider string
	PlaylistInfo
}

type PaginatedPlaylists struct {
	Pagination
	Playlists []PlaylistInfo
}
type EnrichedPaginatedPlaylists struct {
	Provider string
	PaginatedPlaylists
	Playlists []EnrichedPlaylistInfo
}

// -- SearchSetup -- //

type SearchSetup struct {
	Songs     PaginatedSongs
	Albums    PaginatedAlbums
	Artists   PaginatedArtists
	Playlists PaginatedPlaylists
}

type EnrichedSearchSetup struct {
	Songs     EnrichedPaginatedSongs
	Albums    EnrichedPaginatedAlbums
	Artists   EnrichedPaginatedArtists
	Playlists EnrichedPaginatedPlaylists
}

type SearchSetupResult struct {
	ItemFound      TypedItem
	ProviderResult map[string]EnrichedSearchSetup
}

// -- URL -- //

type TypedItem struct {
	Type DataType
	Data any
}

type URLItem struct {
	Type DataType
	ID   string
}
