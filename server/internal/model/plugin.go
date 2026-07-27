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

type SongInfo struct {
	ID              string
	Title           string
	Duration        uint
	ReplayGain      float64
	Peak            float64
	AlbumReplayGain float64
	AlbumPeak       float64
	ReleaseDate     time.Time
	TrackNumber     uint
	VolumeNumber    uint
	AudioQuality    Quality
	Popularity      uint
	Explicit        bool
	Isrc            string
	Artists         []ArtistInfo
	Album           AlbumInfo
}

type EnrichedSongInfo struct {
	Provider string
	SongInfo
}

type AlbumInfo struct {
	ID            string
	Title         string
	Duration      uint
	ReleaseDate   time.Time
	NumberTracks  uint
	NumberVolumes uint
	CoverUrl      string
	AudioQuality  Quality
	Popularity    uint
	Explicit      bool
	Artists       []ArtistInfo
}
type EnrichedAlbumInfo struct {
	Provider string
	AlbumInfo
}

type AlbumSongs struct {
	Songs []SongInfo
}
type EnrichedAlbumSongs struct {
	Provider string
	AlbumSongs
	Songs []EnrichedSongInfo
}

type ArtistInfo struct {
	ID         string
	Name       string
	PictureUrl string
	Popularity uint
}

type EnrichedArtistInfo struct {
	Provider string
	Followed uuid.UUID
	ArtistInfo
}

type ArtistAlbums struct {
	Albums  []AlbumInfo
	Ep      []AlbumInfo
	Singles []AlbumInfo
}

type EnrichedArtistAlbums struct {
	Provider string
	ArtistAlbums
	Albums  []EnrichedAlbumInfo
	Ep      []EnrichedAlbumInfo
	Singles []EnrichedAlbumInfo
}

type PlaylistInfo struct {
	ID             string
	Title          string
	Description    string
	Duration       uint
	LastUpdated    time.Time
	NumberOfTracks uint
	CoverURL       string
	Popularity     uint
}

type EnrichedPlaylistInfo struct {
	Provider string
	PlaylistInfo
}

type PlaylistSongs struct {
	Songs []SongInfo
}
type EnrichedPlaylistSongs struct {
	Provider string
	PlaylistSongs
	Songs []EnrichedSongInfo
}

type TypedItem struct {
	Type DataType
	Data any
}

type Search struct {
	Songs     []SongInfo
	Albums    []AlbumInfo
	Artists   []ArtistInfo
	Playlists []PlaylistInfo
}

type EnrichedSearch struct {
	Songs     []EnrichedSongInfo
	Albums    []EnrichedAlbumInfo
	Artists   []EnrichedArtistInfo
	Playlists []EnrichedPlaylistInfo
}

type SearchResult struct {
	ItemFound      TypedItem
	ProviderResult map[string]EnrichedSearch
}

type UrlItem struct {
	Type DataType
	ID   string
}
