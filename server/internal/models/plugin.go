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

type SongData struct {
	Id                  string           `json:"id"`
	Title               string           `json:"title"`
	Duration            uint             `json:"duration"`
	ReleaseDate         string           `json:"releaseDate"`
	TrackNumber         uint             `json:"trackNumber"`
	VolumeNumber        uint             `json:"volumeNumber"`
	MaximalAudioQuality string           `json:"maximalAudioQuality"`
	Popularity          uint             `json:"popularity"`
	Isrc                string           `json:"isrc"`
	CoverUrl            string           `json:"coverUrl"`
	Artists             []SongDataArtist `json:"artists"`
	Album               SongDataAlbum    `json:"album"`
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
	Id                  string            `json:"id"`
	Title               string            `json:"title"`
	Duration            uint              `json:"duration"`
	ReleaseDate         string            `json:"releaseDate"`
	NumberTracks        uint              `json:"numberTracks"`
	NumberVolumes       uint              `json:"numberVolumes"`
	CoverUrl            string            `json:"coverUrl"`
	MaximalAudioQuality string            `json:"maximalAudioQuality"`
	Artists             []AlbumDataArtist `json:"artists"`
	Songs               []AlbumDataSong   `json:"songs"`
}

type AlbumDataArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AlbumDataSong struct {
	Id                  string           `json:"id"`
	Title               string           `json:"title"`
	Duration            uint             `json:"duration"`
	TrackNumber         uint             `json:"trackNumber"`
	VolumeNumber        uint             `json:"volumeNumber"`
	MaximalAudioQuality string           `json:"maximalAudioQuality"`
	Artists             []SongDataArtist `json:"artists"`
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
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Duration    uint              `json:"duration"`
	ReleaseDate string            `json:"releaseDate"`
	CoverUrl    string            `json:"coverUrl"`
	Artists     []AlbumDataArtist `json:"artists"`
}

type SearchData struct {
	Songs   []SearchDataSong   `json:"songs"`
	Albums  []SearchDataAlbum  `json:"albums"`
	Artists []SearchDataArtist `json:"artists"`
}

type SearchDataSong struct {
	Id                  string           `json:"id"`
	Title               string           `json:"title"`
	Duration            uint             `json:"duration"`
	MaximalAudioQuality string           `json:"maximalAudioQuality"`
	Popularity          uint             `json:"popularity"`
	Artists             []SongDataArtist `json:"artists"`
	Album               SongDataAlbum    `json:"album"`
}

type SearchDataAlbum struct {
	Id                  string            `json:"id"`
	Title               string            `json:"title"`
	Duration            uint              `json:"duration"`
	CoverUrl            string            `json:"coverUrl"`
	MaximalAudioQuality string            `json:"maximalAudioQuality"`
	Popularity          uint              `json:"popularity"`
	Artists             []AlbumDataArtist `json:"artists"`
}

type SearchDataArtist struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PictureUrl string `json:"pictureUrl"`
	Popularity uint   `json:"popularity"`
}

type Status string

const (
	StatusPending Status = "pending"
	StatusRunning Status = "running"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
	StatusCancel  Status = "cancel"
)

type DownloadData struct {
	Id     uint     `json:"id"`
	Data   SongData `json:"data"`
	Api    string   `json:"api"`
	Status Status   `json:"status"`
}
