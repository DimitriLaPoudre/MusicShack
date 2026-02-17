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
		return nil, fmt.Errorf("metadata.getCover: %w", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("metadata.getCover: io.ReadAll: %w", err)
	}

	return bytes.NewReader(buf), nil
}

func WriteTags(path string, tags map[string][]string, trunc bool) error {
	var writeOption taglib.WriteOption
	if trunc {
		writeOption = taglib.Clear
	} else {
		writeOption = 0
	}
	if err := taglib.WriteTags(path, tags, writeOption); err != nil {
		return fmt.Errorf("metadata.WriteTags: taglib.WriteTags: %w", err)
	} else {
		return nil
	}
}

func WriteCover(path string, reader io.Reader) error {
	img, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("metadata.WriteCover: io.ReadAll: %w", err)
	}
	if err := taglib.WriteImage(path, img); err != nil {
		return fmt.Errorf("metadata.WriteCover: taglib.WriteImage: %w", err)
	} else {
		return nil
	}
}

func FormatMetadata(ctx context.Context, userId uint, path string, data models.SongData) error {
	album, err := plugins.GetAlbum(ctx, userId, data.Provider, data.Album.Id)
	if err != nil {
		return fmt.Errorf("metadata.FormatMetadata: %w", err)
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

	tags := map[string][]string{
		models.TagTitle:        {data.Title},
		models.TagAlbum:        {data.Album.Title},
		models.TagAlbumArtists: albumArtists,
		models.TagArtists:      artists,
		models.TagTrackNumber:  {trackNumber},
		models.TagVolumeNumber: {volumeNumber},
		models.TagReleaseDate:  {album.ReleaseDate},
		models.TagExplicit:     {explicit},
		models.TagAlbumGain:    {albumGain},
		models.TagAlbumPeak:    {albumPeak},
		models.TagTrackGain:    {trackGain},
		models.TagTrackPeak:    {trackPeak},
		models.TagISRC:         {data.Isrc},
	}

	if err := WriteTags(path, tags, false); err != nil {
		return fmt.Errorf("metadata.FormatMetadata: %w", err)
	}

	img, err := getCover(ctx, data.Album.CoverUrl)
	if err != nil {
		return fmt.Errorf("metadata.FormatMetadata: %w", err)
	}

	if err := WriteCover(path, img); err != nil {
		return fmt.Errorf("metadata.FormatMetadata: %w", err)
	}
	return nil
}

func ReadTags(path string) (map[string][]string, error) {
	if tags, err := taglib.ReadTags(path); err != nil {
		return tags, fmt.Errorf("metadata.ReadTags: %w", err)
	} else {
		return tags, nil
	}
}

func ReadCover(path string) ([]byte, error) {
	img, err := taglib.ReadImage(path)
	if err != nil {
		return nil, fmt.Errorf("metadata.ReadCover: %w", err)
	} else {
		return img, nil
	}
}
