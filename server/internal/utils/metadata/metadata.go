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

func getCover(ctx context.Context, url string) (*[]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("getCover: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
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

func FormatMetadata(ctx context.Context, userId uint, path string, data models.SongData) error {
	album, err := plugins.GetAlbum(ctx, userId, data.Provider, data.Album.Id)
	if err != nil {
		return fmt.Errorf("FormatMetadata: api.Album: %w", err)
	}
	var albumArtists []string
	for _, artist := range album.Artists {
		albumArtists = append(albumArtists, artist.Name)
	}

	var artists []string
	for _, artist := range data.Artists {
		artists = append(artists, artist.Name)
	}

	explicit := "false"
	if data.Explicit {
		explicit = "true"
	}
	trackNumber := strconv.FormatUint(uint64(data.TrackNumber), 10)
	volumeNumber := strconv.FormatUint(uint64(data.VolumeNumber), 10)
	trackGain := strconv.FormatFloat(data.ReplayGain, 'f', -1, 64)
	trackPeak := strconv.FormatFloat(data.Peak, 'f', -1, 64)
	albumGain := strconv.FormatFloat(data.AlbumReplayGain, 'f', -1, 64)
	albumPeak := strconv.FormatFloat(data.AlbumPeak, 'f', -1, 64)

	if err := taglib.WriteTags(path, map[string][]string{
		taglib.Title:            {data.Title},
		taglib.Artists:          artists,
		"ALBUMARTISTS":          albumArtists,
		taglib.Album:            {data.Album.Title},
		taglib.TrackNumber:      {trackNumber},
		taglib.DiscNumber:       {volumeNumber},
		taglib.ReleaseDate:      {album.ReleaseDate},
		"itunesadvisory":        {explicit},
		"replaygain_album_gain": {albumGain},
		"replaygain_album_peak": {albumPeak},
		"replaygain_track_gain": {trackGain},
		"replaygain_track_peak": {trackPeak},
		taglib.ISRC:             {data.Isrc},
	}, taglib.Clear); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteTags: %w", err)
	}

	image, err := getCover(ctx, data.Album.CoverUrl)
	if err != nil {
		return fmt.Errorf("FormatMetadata: getCover: %w", err)
	}

	if err := taglib.WriteImage(path, *image); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteImage: %w", err)
	}
	return nil
}
