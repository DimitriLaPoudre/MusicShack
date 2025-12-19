package models

import (
	"context"
	"io"
)

type Plugin interface {
	Name() string
	Download(context.Context, string, string, chan<- SongData) (io.ReadCloser, string, error)
	Song(context.Context, string) (SongData, error)
	Album(context.Context, string) (AlbumData, error)
	Artist(context.Context, string) (ArtistData, error)
	Search(context.Context, string, string, string) (SearchData, error)
	Cover(context.Context, string) (string, error)
	Lyrics(context.Context, string) (string, string, error)
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
	Id           string
	Title        string
	Duration     uint
	ReleaseDate  string
	TrackNumber  uint
	VolumeNumber uint
	AudioQuality string
	Artist       struct {
		Id   string
		Name string
	}
	Artists []struct {
		Id   string
		Name string
	}
	Album struct {
		Id       string
		Title    string
		CoverUrl string
	}
	DownloadUrl string
}

type AlbumData struct {
	Id            string
	Title         string
	Duration      uint
	ReleaseDate   string
	NumberTracks  uint
	NumberVolumes uint
	Type          string
	CoverUrl      string
	AudioQuality  string
	Artist        struct {
		Id   string
		Name string
	}
	Artists []struct {
		Id   string
		Name string
	}
	Limit       uint
	Offset      uint
	NumberSongs uint
	Songs       []struct {
		Id           string
		Title        string
		Duration     uint
		TrackNumber  uint
		VolumeNumber uint
		Artists      []struct {
			Id   string
			Name string
		}
	}
}

type ArtistData struct {
	Id         string
	Name       string
	PictureUrl string
	Albums     []AlbumData
	// Albums     []struct {
	// 	Id       string
	// 	Title    string
	// 	CoverUrl string
	// 	Artists  []struct {
	// 		Id   string
	// 		Name string
	// 	}
	// }
	Ep []struct {
		Id       string
		Title    string
		CoverUrl string
	}
	Singles []struct {
		Id       string
		Title    string
		CoverUrl string
	}
}

type SearchData struct {
	Songs []struct {
		Id       string
		Title    string
		CoverUrl string
		Artists  []struct {
			Id   string
			Name string
		}
	}
	Albums []struct {
		Id       string
		Title    string
		CoverUrl string
		Artists  []struct {
			Id   string
			Name string
		}
	}
	Artists []struct {
		Id         string
		Name       string
		PictureUrl string
	}
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
