package metadata

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"go.senan.xyz/taglib"
)

func getCover(ctx context.Context, url string) (*[]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := utils.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("getCover: %w", err)
	}
	defer resp.Body.Close()

	image, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getCover: io.ReadAll: %w", err)
	}
	return &image, nil
}

func ApplyMetadata(path string, info models.MetadataInfo) error {
	if err := taglib.WriteTags(path, map[string][]string{
		taglib.Title:            {info.Title},
		taglib.Artists:          info.Artists,
		"ALBUMARTISTS":          info.AlbumArtists,
		taglib.Album:            {info.Album},
		taglib.TrackNumber:      {info.TrackNumber},
		taglib.DiscNumber:       {info.VolumeNumber},
		taglib.ReleaseDate:      {info.ReleaseDate},
		"itunesadvisory":        {info.Explicit},
		"replaygain_album_gain": {info.AlbumGain},
		"replaygain_album_peak": {info.AlbumPeak},
		"replaygain_track_gain": {info.TrackGain},
		"replaygain_track_peak": {info.TrackPeak},
		taglib.ISRC:             {info.Isrc},
	}, taglib.Clear); err != nil {
		return fmt.Errorf("ApplyMetadata: taglib.WriteTags: %w", err)
	} else {
		return nil
	}
}

func ApplyCover(path string, img *[]byte) error {
	if err := taglib.WriteImage(path, *img); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteImage: %w", err)
	} else {
		return nil
	}
}

func FormatMetadata(ctx context.Context, userId uint, path string, data models.SongData) error {
	album, err := plugins.GetAlbum(ctx, userId, data.Provider, data.Album.Id)
	if err != nil {
		return fmt.Errorf("FormatMetadata: %w", err)
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
		return fmt.Errorf("FormatMetadata: %w", err)
	}

	if err := taglib.WriteImage(path, *image); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteImage: %w", err)
	}
	return nil
}
