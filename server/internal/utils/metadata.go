package utils

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"go.senan.xyz/taglib"
)

func getCover(data *models.SongData) (*[]byte, error) {
	resp, err := http.Get(data.Album.CoverUrl)
	if err != nil {
		return nil, fmt.Errorf("getCover: %w", err)
	}
	defer resp.Body.Close()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getCover: %w", err)
	}

	return &image, nil
}

func FormatMetadata(path string, data models.SongData) error {
	var artists []string
	for _, artist := range data.Artists {
		artists = append(artists, artist.Name)
	}
	trackNumber := strconv.FormatUint(uint64(data.TrackNumber), 10)
	volumeNumber := strconv.FormatUint(uint64(data.VolumeNumber), 10)

	if err := taglib.WriteTags(path, map[string][]string{
		taglib.Title:       {data.Title},
		taglib.Artist:      {data.Artists[0].Name},
		taglib.Artists:     artists,
		taglib.Album:       {data.Album.Title},
		taglib.AlbumArtist: {data.Artists[0].Name},
		taglib.TrackNumber: {trackNumber},
		taglib.DiscNumber:  {volumeNumber},
		taglib.Date:        {data.ReleaseDate},
	}, taglib.Clear); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteTags: %w", err)
	}

	image, err := getCover(&data)
	if err != nil {
		return fmt.Errorf("FormatMetadata: getCover: %w", err)
	}

	if err := taglib.WriteImage(path, *image); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteImage: %w", err)
	}
	return nil
}
