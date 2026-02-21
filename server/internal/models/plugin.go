package models

import (
	"context"
	"io"
)

type Plugin interface {
	Name() string
	Provider() string
	Priority() int
	Status(ctx context.Context, url string) error
	Download(context.Context, uint, string) (io.ReadCloser, string, error)
	Song(context.Context, uint, string) (SongData, error)
	Playlist(context.Context, uint, string) (PlaylistData, error)
	Album(context.Context, uint, string) (AlbumData, error)
	Artist(context.Context, uint, string) (ArtistData, error)
	Search(context.Context, uint, string, string, string) (SearchData, error)
	Url(context.Context, uint, string) (UrlItem, error)
	Lyrics(context.Context, uint, string) (string, string, error)
}

type Quality struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type SongData struct {
	Provider        string           `json:"provider"`
	Api             string           `json:"api"`
	Downloaded      bool             `json:"downloaded"`
	Id              string           `json:"id"`
	Title           string           `json:"title"`
	Duration        uint             `json:"duration"`
	ReplayGain      float64          `json:"replayGain"`
	Peak            float64          `json:"peak"`
	AlbumReplayGain float64          `json:"albumReplayGain"`
	AlbumPeak       float64          `json:"albumPeak"`
	ReleaseDate     string           `json:"releaseDate"`
	TrackNumber     uint             `json:"trackNumber"`
	VolumeNumber    uint             `json:"volumeNumber"`
	AudioQuality    Quality          `json:"audioQuality"`
	Explicit        bool             `json:"explicit"`
	Popularity      uint             `json:"popularity"`
	Isrc            string           `json:"isrc"`
	Artists         []SongDataArtist `json:"artists"`
	Album           SongDataAlbum    `json:"album"`
}

type SongDataArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SongDataAlbum struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

type PlaylistData struct {
	Provider       string             `json:"provider"`
	Api            string             `json:"api"`
	Downloaded     bool               `json:"downloaded"`
	Id             string             `json:"id"`
	Title          string             `json:"title"`
	Description    string             `json:"description"`
	Duration       uint               `json:"duration"`
	LastUpdated    string             `json:"lastUpdated"`
	NumberOfTracks uint               `json:"numberOfTracks"`
	CoverURL       string             `json:"coverUrl"`
	Songs          []PlaylistDataSong `json:"songs"`
}

type PlaylistDataSong struct {
	Downloaded   bool             `json:"downloaded"`
	Id           string           `json:"id"`
	Title        string           `json:"title"`
	Duration     uint             `json:"duration"`
	AudioQuality Quality          `json:"audioQuality"`
	Explicit     bool             `json:"explicit"`
	Isrc         string           `json:"isrc"`
	Artists      []SongDataArtist `json:"artists"`
}

type AlbumData struct {
	Provider      string            `json:"provider"`
	Api           string            `json:"api"`
	Downloaded    bool              `json:"downloaded"`
	Id            string            `json:"id"`
	Title         string            `json:"title"`
	Duration      uint              `json:"duration"`
	ReleaseDate   string            `json:"releaseDate"`
	NumberTracks  uint              `json:"numberTracks"`
	NumberVolumes uint              `json:"numberVolumes"`
	CoverUrl      string            `json:"coverUrl"`
	AudioQuality  Quality           `json:"audioQuality"`
	Explicit      bool              `json:"explicit"`
	Artists       []AlbumDataArtist `json:"artists"`
	Songs         []AlbumDataSong   `json:"songs"`
}

type AlbumDataArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AlbumDataSong struct {
	Downloaded   bool             `json:"downloaded"`
	Id           string           `json:"id"`
	Title        string           `json:"title"`
	Duration     uint             `json:"duration"`
	TrackNumber  uint             `json:"trackNumber"`
	VolumeNumber uint             `json:"volumeNumber"`
	AudioQuality Quality          `json:"audioQuality"`
	Explicit     bool             `json:"explicit"`
	Isrc         string           `json:"isrc"`
	Artists      []SongDataArtist `json:"artists"`
}

type ArtistData struct {
	Provider   string            `json:"provider"`
	Api        string            `json:"api"`
	Followed   uint              `json:"followed"`
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	PictureUrl string            `json:"pictureUrl"`
	Albums     []ArtistDataAlbum `json:"albums"`
	Ep         []ArtistDataAlbum `json:"ep"`
	Singles    []ArtistDataAlbum `json:"singles"`
}

type ArtistDataAlbum struct {
	Downloaded   bool              `json:"downloaded"`
	Id           string            `json:"id"`
	Title        string            `json:"title"`
	Duration     uint              `json:"duration"`
	ReleaseDate  string            `json:"releaseDate"`
	CoverUrl     string            `json:"coverUrl"`
	AudioQuality Quality           `json:"audioQuality"`
	Explicit     bool              `json:"explicit"`
	Artists      []AlbumDataArtist `json:"artists"`
}

type SearchData struct {
	Url       UrlItem
	Songs     []SearchDataSong     `json:"songs"`
	Albums    []SearchDataAlbum    `json:"albums"`
	Artists   []SearchDataArtist   `json:"artists"`
	Playlists []SearchDataPlaylist `json:"playlists"`
}

type SearchDataSong struct {
	Downloaded   bool             `json:"downloaded"`
	Id           string           `json:"id"`
	Title        string           `json:"title"`
	Duration     uint             `json:"duration"`
	AudioQuality Quality          `json:"audioQuality"`
	Popularity   uint             `json:"popularity"`
	Explicit     bool             `json:"explicit"`
	Isrc         string           `json:"isrc"`
	Artists      []SongDataArtist `json:"artists"`
	Album        SongDataAlbum    `json:"album"`
}

type SearchDataAlbum struct {
	Downloaded   bool              `json:"downloaded"`
	Id           string            `json:"id"`
	Title        string            `json:"title"`
	Duration     uint              `json:"duration"`
	CoverUrl     string            `json:"coverUrl"`
	AudioQuality Quality           `json:"audioQuality"`
	Explicit     bool              `json:"explicit"`
	Popularity   uint              `json:"popularity"`
	Artists      []AlbumDataArtist `json:"artists"`
}

type SearchDataArtist struct {
	Followed   uint   `json:"followed"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
	Popularity uint   `json:"popularity"`
}

type SearchDataPlaylist struct {
	Downloaded bool   `json:"downloaded"`
	ID         string `json:"id"`
	Title      string `json:"title"`
	Duration   uint   `json:"duration"`
	CoverURL   string `json:"coverUrl"`
}

type Type string

const (
	TypeSong   Type = "song"
	TypeAlbum  Type = "album"
	TypeArtist Type = "artist"
)

type UrlItem struct {
	Provider string `json:"provider"`
	Type     Type   `json:"type"`
	Id       string `json:"id"`
}
