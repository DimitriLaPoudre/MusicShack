package utils

import (
	"fmt"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"go.senan.xyz/taglib"
)

func FormatMetadata(path string, data models.SongData) error {
	var artists []string
	for _, artist := range data.Artists {
		artists = append(artists, artist.Name)
	}
	trackNumber := strconv.FormatUint(uint64(data.TrackNumber), 10)
	volumeNumber := strconv.FormatUint(uint64(data.VolumeNumber), 10)

	if err := taglib.WriteTags(path, map[string][]string{
		taglib.Title:       {data.Title},
		taglib.Artist:      {data.Artist.Name},
		taglib.Artists:     artists,
		taglib.Album:       {data.Album.Title},
		taglib.AlbumArtist: {data.Artist.Name},
		taglib.TrackNumber: {trackNumber},
		taglib.DiscNumber:  {volumeNumber},
		taglib.Date:        {data.ReleaseDate},
	}, taglib.Clear); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteTags: %w", err)
	}
	return nil
}
