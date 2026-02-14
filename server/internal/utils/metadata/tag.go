package metadata

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

const (
	Title        string = "TITLE"
	Artists      string = "ARTISTS"
	AlbumArtists string = "ALBUMARTISTS"
	Album        string = "ALBUM"
	TrackNumber  string = "TRACKNUMBER"
	VolumeNumber string = "DISCNUMBER"
	ReleaseDate  string = "RELEASEDATE"
	Explicit     string = "ITUNESADVISORY"
	AlbumGain    string = "REPLAYGAIN_ALBUM_GAIN"
	AlbumPeak    string = "REPLAYGAIN_ALBUM_PEAK"
	TrackGain    string = "REPLAYGAIN_TRACK_GAIN"
	TrackPeak    string = "REPLAYGAIN_TRACK_PEAK"
	ISRC         string = "ISRC"
)

func MetadataInfoToTags(info models.MetadataInfo) map[string][]string {
	return map[string][]string{
		Title:        {info.Title},
		Album:        {info.Title},
		AlbumArtists: info.AlbumArtists,
		Artists:      info.Artists,
		TrackNumber:  {info.TrackNumber},
		VolumeNumber: {info.VolumeNumber},
		ReleaseDate:  {info.ReleaseDate},
		Explicit:     {info.Explicit},
		AlbumGain:    {info.AlbumGain},
		AlbumPeak:    {info.AlbumPeak},
		TrackGain:    {info.TrackGain},
		TrackPeak:    {info.TrackPeak},
		ISRC:         {info.Isrc},
	}
}
