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

// "tags": {
//             "title": "CASINO (feat. La Fève)",
//             "artist": "thaHomey; Skuna; La Fève",
//             "album_artist": "thaHomey",
//             "album": "HEAT",
//             "track": "4",
//             "disc": "1",
//             "date": "2022",
//             "year": "2022",
//             "ISRC": "FR9W12220849",
//             "REPLAYGAIN_TRACK_GAIN": "-7.77 dB",
//             "REPLAYGAIN_TRACK_PEAK": "0.987245",
//             "REPLAYGAIN_ALBUM_GAIN": "-8.98 dB",
//             "REPLAYGAIN_ALBUM_PEAK": "0.987245",
//             "comment": "Downloaded from music.binimum.org/tidal.squid.wtf",
//             "encoder": "Lavf59.27.100"
//         }

type SongData struct {
	Id                  string
	Title               string
	Duration            uint
	ReleaseDate         string
	TrackNumber         uint
	VolumeNumber        uint
	MaximalAudioQuality string
	Popularity          uint
	Isrc                string
	CoverUrl            string
	Artists             []SongDataArtist
	Album               SongDataAlbum
}

type SongDataArtist struct {
	Id   string
	Name string
}

type SongDataAlbum struct {
	Id       string
	Title    string
	CoverUrl string
}

type AlbumData struct {
	Id                  string
	Title               string
	Duration            uint
	ReleaseDate         string
	NumberTracks        uint
	NumberVolumes       uint
	CoverUrl            string
	MaximalAudioQuality string
	Artists             []AlbumDataArtist
	Songs               []AlbumDataSong
}

type AlbumDataArtist struct {
	Id   string
	Name string
}

type AlbumDataSong struct {
	Id                  string
	Title               string
	Duration            uint
	TrackNumber         uint
	VolumeNumber        uint
	MaximalAudioQuality string
	Artists             []SongDataArtist
}

type ArtistData struct {
	Id         string
	Name       string
	PictureUrl string
	Albums     []ArtistDataAlbum
	Ep         []ArtistDataAlbum
	Singles    []ArtistDataAlbum
}

type ArtistDataAlbum struct {
	Id          string
	Title       string
	Duration    uint
	ReleaseDate string
	CoverUrl    string
	Artists     []AlbumDataArtist
}

type SearchData struct {
	Songs   []SearchDataSong
	Albums  []SearchDataAlbum
	Artists []SearchDataArtist
}

type SearchDataSong struct {
	Id                  string
	Title               string
	Duration            uint
	MaximalAudioQuality string
	Popularity          uint
	Artists             []SongDataArtist
	Album               SongDataAlbum
}

type SearchDataAlbum struct {
	Id                  string
	Title               string
	Duration            uint
	CoverUrl            string
	MaximalAudioQuality string
	Popularity          uint
	Artists             []AlbumDataArtist
}

type SearchDataArtist struct {
	Id         string
	Name       string
	PictureUrl string
	Popularity uint
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
	Id     uint
	Data   SongData
	Api    string
	Status Status
}
