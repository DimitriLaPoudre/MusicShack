package metadata

import (
	"bytes"
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

func getCover(ctx context.Context, url string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := utils.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("getCover: %w", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("getCover: io.ReadAll: %w", err)
	}

	return bytes.NewReader(buf), nil
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

func ApplyCover(path string, reader io.Reader) error {
	img, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("ApplyCover: io.ReadAll: %w", err)
	}
	if err := taglib.WriteImage(path, img); err != nil {
		return fmt.Errorf("ApplyCover: taglib.WriteImage: %w", err)
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

	info := models.MetadataInfo{
		Title:        data.Title,
		Album:        data.Album.Title,
		AlbumArtists: albumArtists,
		Artists:      artists,
		TrackNumber:  trackNumber,
		VolumeNumber: volumeNumber,
		ReleaseDate:  album.ReleaseDate,
		Explicit:     explicit,
		AlbumGain:    albumGain,
		AlbumPeak:    albumPeak,
		TrackGain:    trackGain,
		TrackPeak:    trackPeak,
		Isrc:         data.Isrc,
	}

	if err := ApplyMetadata(path, info); err != nil {
		return fmt.Errorf("FormatMetadata: %w", err)
	}

	img, err := getCover(ctx, data.Album.CoverUrl)
	if err != nil {
		return fmt.Errorf("FormatMetadata: %w", err)
	}

	if err := ApplyCover(path, img); err != nil {
		return fmt.Errorf("FormatMetadata: taglib.WriteImage: %w", err)
	}
	return nil
}
