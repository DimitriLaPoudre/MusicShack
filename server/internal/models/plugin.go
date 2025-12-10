package models

import "context"

type Plugin interface {
	Name() string
	DownloadSong(context.Context, uint, string, string) error
	DownloadAlbum(context.Context, uint, string, string) error
	DownloadArtist(context.Context, uint, string, string) error
	Song(context.Context, string) (SongData, error)
	Album(context.Context, string) (AlbumData, error)
	Artist(context.Context, string) (ArtistData, error)
	Search(context.Context, string, string, string) (SearchData, error)
	Cover(context.Context, string) (string, error)
	Lyrics(context.Context, string) (string, string, error)
}

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
	Albums     []struct {
		Id       string
		Title    string
		CoverUrl string
		Artists  []struct {
			Id   string
			Name string
		}
	}
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
