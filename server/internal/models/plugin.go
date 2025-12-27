package models

import (
	"context"
	"io"
)

type Plugin interface {
	Name() string
	Download(context.Context, uint, string, string, chan<- SongData) (io.ReadCloser, string, error)
	Song(context.Context, uint, string) (SongData, error)
	Album(context.Context, uint, string) (AlbumData, error)
	Artist(context.Context, uint, string) (ArtistData, error)
	Search(context.Context, uint, string, string, string) (SearchData, error)
	Lyrics(context.Context, uint, string) (string, string, error)
}

type Quality uint

const (
	QualityLow Quality = 1
	QualityHigh
	QualityLossless
	QualityHiresLossless
)

type SongData struct {
	Id           string           `json:"id"`
	Title        string           `json:"title"`
	Duration     uint             `json:"duration"`
	ReleaseDate  string           `json:"releaseDate"`
	TrackNumber  uint             `json:"trackNumber"`
	VolumeNumber uint             `json:"volumeNumber"`
	AudioQuality Quality          `json:"audioQuality"`
	Popularity   uint             `json:"popularity"`
	Isrc         string           `json:"isrc"`
	CoverUrl     string           `json:"coverUrl"`
	Artists      []SongDataArtist `json:"artists"`
	Album        SongDataAlbum    `json:"album"`
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

type AlbumData struct {
	Id            string            `json:"id"`
	Title         string            `json:"title"`
	Duration      uint              `json:"duration"`
	ReleaseDate   string            `json:"releaseDate"`
	NumberTracks  uint              `json:"numberTracks"`
	NumberVolumes uint              `json:"numberVolumes"`
	CoverUrl      string            `json:"coverUrl"`
	AudioQuality  Quality           `json:"audioQuality"`
	Artists       []AlbumDataArtist `json:"artists"`
	Songs         []AlbumDataSong   `json:"songs"`
}

type AlbumDataArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AlbumDataSong struct {
	Id           string           `json:"id"`
	Title        string           `json:"title"`
	Duration     uint             `json:"duration"`
	TrackNumber  uint             `json:"trackNumber"`
	VolumeNumber uint             `json:"volumeNumber"`
	AudioQuality Quality          `json:"audioQuality"`
	Artists      []SongDataArtist `json:"artists"`
}

type ArtistData struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	PictureUrl string            `json:"pictureUrl"`
	Albums     []ArtistDataAlbum `json:"albums"`
	Ep         []ArtistDataAlbum `json:"ep"`
	Singles    []ArtistDataAlbum `json:"singles"`
}

type ArtistDataAlbum struct {
	Id           string            `json:"id"`
	Title        string            `json:"title"`
	Duration     uint              `json:"duration"`
	ReleaseDate  string            `json:"releaseDate"`
	CoverUrl     string            `json:"coverUrl"`
	AudioQuality Quality           `json:"audioQuality"`
	Artists      []AlbumDataArtist `json:"artists"`
}

type SearchData struct {
	Songs   []SearchDataSong   `json:"songs"`
	Albums  []SearchDataAlbum  `json:"albums"`
	Artists []SearchDataArtist `json:"artists"`
}

type SearchDataSong struct {
	Id           string           `json:"id"`
	Title        string           `json:"title"`
	Duration     uint             `json:"duration"`
	AudioQuality Quality          `json:"audioQuality"`
	Popularity   uint             `json:"popularity"`
	Artists      []SongDataArtist `json:"artists"`
	Album        SongDataAlbum    `json:"album"`
}

type SearchDataAlbum struct {
	Id           string            `json:"id"`
	Title        string            `json:"title"`
	Duration     uint              `json:"duration"`
	CoverUrl     string            `json:"coverUrl"`
	AudioQuality Quality           `json:"audioQuality"`
	Popularity   uint              `json:"popularity"`
	Artists      []AlbumDataArtist `json:"artists"`
}

type SearchDataArtist struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
	Popularity uint   `json:"popularity"`
}
