package response

import "github.com/Ascension-EIP/Ascension/apps/server/internal/model"

type Quality struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Song struct {
	Provider        string       `json:"provider"`
	Api             string       `json:"api"`
	Downloaded      bool         `json:"downloaded"`
	Id              string       `json:"id"`
	Title           string       `json:"title"`
	Duration        uint         `json:"duration"`
	ReplayGain      float64      `json:"replayGain"`
	Peak            float64      `json:"peak"`
	AlbumReplayGain float64      `json:"albumReplayGain"`
	AlbumPeak       float64      `json:"albumPeak"`
	ReleaseDate     string       `json:"releaseDate"`
	TrackNumber     uint         `json:"trackNumber"`
	VolumeNumber    uint         `json:"volumeNumber"`
	AudioQuality    Quality      `json:"audioQuality"`
	Explicit        bool         `json:"explicit"`
	Popularity      uint         `json:"popularity"`
	Isrc            string       `json:"isrc"`
	Artists         []SongArtist `json:"artists"`
	Album           SongAlbum    `json:"album"`
}

type SongArtist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SongAlbum struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	CoverUrl string `json:"coverUrl"`
}

func SongToResponse(song model.Song) Song {
	artists := []SongArtist{}
	for _, a := range song.Artists {
		artists = append(artists, SongArtist{
			Id:   a.Id,
			Name: a.Name,
		})
	}

	return Song{
		Provider:        song.Provider,
		Api:             song.Api,
		Downloaded:      song.Downloaded,
		Id:              song.Id,
		Title:           song.Title,
		Duration:        song.Duration,
		ReplayGain:      song.ReplayGain,
		Peak:            song.Peak,
		AlbumReplayGain: song.AlbumReplayGain,
		AlbumPeak:       song.AlbumPeak,
		ReleaseDate:     song.ReleaseDate,
		TrackNumber:     song.TrackNumber,
		VolumeNumber:    song.VolumeNumber,
		AudioQuality: Quality{
			Name:  song.AudioQuality.Name,
			Color: song.AudioQuality.Color,
		},
		Explicit:   song.Explicit,
		Popularity: song.Popularity,
		Isrc:       song.Isrc,
		Artists:    artists,
		Album: SongAlbum{
			Id:       song.Album.Id,
			Title:    song.Album.Title,
			CoverUrl: song.Album.CoverUrl,
		},
	}
}
