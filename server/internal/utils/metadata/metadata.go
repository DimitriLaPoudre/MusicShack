package metadata

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
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

func FormatMetadata(userId uint, path string, data models.SongData) error {
	api, ok := plugins.Get(data.Api)
	if !ok {
		return fmt.Errorf("FormatMetadata: plugins.Get: invalid api name")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	album, err := api.Album(ctx, userId, data.Album.Id)
	if err != nil {
		return fmt.Errorf("FormatMetadata: api.Album: %w", err)
	}
	var albumArtists []string
	for _, artist := range album.Artists {
		albumArtists = append(albumArtists, artist.Name)
	}

	albumGain := strconv.FormatFloat(0.0, 'f', -1, 64)
	albumPeak := strconv.FormatFloat(0.0, 'f', -1, 64)

	var artists []string
	for _, artist := range data.Artists {
		artists = append(artists, artist.Name)
	}

	trackNumber := strconv.FormatUint(uint64(data.TrackNumber), 10)
	volumeNumber := strconv.FormatUint(uint64(data.VolumeNumber), 10)
	trackGain := strconv.FormatFloat(data.ReplayGain, 'f', -1, 64)
	trackPeak := strconv.FormatFloat(data.Peak, 'f', -1, 64)

	if err := taglib.WriteTags(path, map[string][]string{
		taglib.Title:            {data.Title},
		taglib.Artists:          artists,
		"ALBUMARTISTS":          albumArtists,
		taglib.Album:            {data.Album.Title},
		taglib.TrackNumber:      {trackNumber},
		taglib.DiscNumber:       {volumeNumber},
		taglib.ReleaseDate:      {data.ReleaseDate},
		"replaygain_album_gain": {albumGain},
		"replaygain_album_peak": {albumPeak},
		"replaygain_track_gain": {trackGain},
		"replaygain_track_peak": {trackPeak},
		taglib.ISRC:             {data.Isrc},
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
