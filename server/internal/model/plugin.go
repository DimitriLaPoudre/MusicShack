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
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Song struct {
	Id              string    `json:"id"`
	Title           string    `json:"title"`
	Duration        uint      `json:"duration"`
	ReplayGain      float64   `json:"replayGain"`
	Peak            float64   `json:"peak"`
	AlbumReplayGain float64   `json:"albumReplayGain"`
	AlbumPeak       float64   `json:"albumPeak"`
	ReleaseDate     time.Time `json:"releaseDate"`
	TrackNumber     uint      `json:"trackNumber"`
	VolumeNumber    uint      `json:"volumeNumber"`
	AudioQuality    Quality   `json:"audioQuality"`
	Popularity      uint      `json:"popularity"`
	Explicit        bool      `json:"explicit"`
	Isrc            string    `json:"isrc"`
	Artists         []Artist  `json:"artists"`
	Album           Album     `json:"album"`
}

type EnrichedSong struct {
	Provider   string `json:"provider"`
	Downloaded bool   `json:"downloaded"`
	Song
}

type Album struct {
	Id            string    `json:"id"`
	Title         string    `json:"title"`
	Duration      uint      `json:"duration"`
	ReleaseDate   time.Time `json:"releaseDate"`
	NumberTracks  uint      `json:"numberTracks"`
	NumberVolumes uint      `json:"numberVolumes"`
	CoverUrl      string    `json:"coverUrl"`
	AudioQuality  Quality   `json:"audioQuality"`
	Popularity    uint      `json:"popularity"`
	Explicit      bool      `json:"explicit"`
	Artists       []Artist  `json:"artists"`
	Songs         []Song    `json:"songs"`
}

type EnrichedAlbum struct {
	Provider   string `json:"provider"`
	Downloaded bool   `json:"downloaded"`
	Album
	Songs []EnrichedSong `json:"songs"`
}

type Artist struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	PictureUrl string  `json:"pictureUrl"`
	Popularity uint    `json:"popularity"`
	Albums     []Album `json:"albums"`
	Ep         []Album `json:"ep"`
	Singles    []Album `json:"singles"`
}

type EnrichedArtist struct {
	Provider string    `json:"provider"`
	Followed uuid.UUID `json:"followed"`
	Artist
	Albums  []EnrichedAlbum `json:"albums"`
	Ep      []EnrichedAlbum `json:"ep"`
	Singles []EnrichedAlbum `json:"singles"`
}

type Playlist struct {
	Id             string    `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Duration       uint      `json:"duration"`
	LastUpdated    time.Time `json:"lastUpdated"`
	NumberOfTracks uint      `json:"numberOfTracks"`
	CoverURL       string    `json:"coverUrl"`
	Popularity     uint      `json:"popularity"`
	Songs          []Song    `json:"songs"`
}

type EnrichedPlaylist struct {
	Provider   string `json:"provider"`
	Downloaded bool   `json:"downloaded"`
	Playlist
	Songs []EnrichedSong `json:"songs"`
}

type TypedItem struct {
	Type DataType
	Data any
}

type Search struct {
	Songs     []Song     `json:"songs"`
	Albums    []Album    `json:"albums"`
	Artists   []Artist   `json:"artists"`
	Playlists []Playlist `json:"playlists"`
}

type EnrichedSearch struct {
	Songs     []EnrichedSong     `json:"songs"`
	Albums    []EnrichedAlbum    `json:"albums"`
	Artists   []EnrichedArtist   `json:"artists"`
	Playlists []EnrichedPlaylist `json:"playlists"`
}

type SearchResult struct {
	ItemFound      TypedItem
	ProviderResult map[string]EnrichedSearch
}

type UrlItem struct {
	Type DataType `json:"type"`
	Id   string   `json:"id"`
}
