package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/network"
	"go.senan.xyz/taglib"
)

type MetadataService struct {
	plugin *PluginService
}

func NewMetadataService(plugin *PluginService) MetadataService {
	return MetadataService{
		plugin: plugin,
	}
}

func (s *MetadataService) getCover(ctx context.Context, url string) (io.Reader, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := network.Fetch(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetch cover: %w", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read cover from response: %w", err)
	}

	return bytes.NewReader(buf), nil
}

func (s *MetadataService) WriteTags(path string, tags map[string][]string, trunc bool) error {
	var writeOption taglib.WriteOption
	if trunc {
		writeOption = taglib.Clear
	} else {
		writeOption = 0
	}

	if err := taglib.WriteTags(path, tags, writeOption); err != nil {
		return err
	}

	return nil
}

func (s *MetadataService) WriteCover(path string, reader io.Reader) error {
	img, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read cover: %w", err)
	}

	if err := taglib.WriteImage(path, img); err != nil {
		return fmt.Errorf("write cover: %w", err)
	}

	return nil
}

func (s *MetadataService) FormatMetadata(ctx context.Context, user model.User, provider string, path string, songInfo model.SongInfo) error {
	var albumArtists []string
	for _, artist := range songInfo.Album.Artists {
		albumArtists = append(albumArtists, artist.Name)
	}

	var artists []string
	for _, artist := range songInfo.Artists {
		artists = append(artists, artist.Name)
	}

	explicit := "0"
	if songInfo.Explicit {
		explicit = "1"
	}
	trackNumber := strconv.FormatUint(uint64(songInfo.TrackNumber), 10)
	volumeNumber := strconv.FormatUint(uint64(songInfo.VolumeNumber), 10)
	trackGain := strconv.FormatFloat(songInfo.ReplayGain, 'f', -1, 64)
	trackPeak := strconv.FormatFloat(songInfo.Peak, 'f', -1, 64)
	albumGain := strconv.FormatFloat(songInfo.AlbumReplayGain, 'f', -1, 64)
	albumPeak := strconv.FormatFloat(songInfo.AlbumPeak, 'f', -1, 64)

	tags := map[string][]string{
		model.TagTitle:        {songInfo.Title},
		model.TagAlbum:        {songInfo.Album.Title},
		model.TagAlbumArtists: albumArtists,
		model.TagArtists:      artists,
		model.TagTrackNumber:  {trackNumber},
		model.TagVolumeNumber: {volumeNumber},
		model.TagReleaseDate:  {songInfo.Album.ReleaseDate.String()},
		model.TagExplicit:     {explicit},
		model.TagAlbumGain:    {albumGain},
		model.TagAlbumPeak:    {albumPeak},
		model.TagTrackGain:    {trackGain},
		model.TagTrackPeak:    {trackPeak},
		model.TagISRC:         {songInfo.Isrc},
	}

	if err := s.WriteTags(path, tags, false); err != nil {
		return fmt.Errorf("save song tags: %w", err)
	}

	img, err := s.getCover(ctx, songInfo.Album.CoverUrl)
	if err != nil {
		return fmt.Errorf("get song cover: %w", err)
	}

	if err := s.WriteCover(path, img); err != nil {
		return fmt.Errorf("save song cover: %w", err)
	}
	return nil
}

func (s *MetadataService) ReadTags(path string) (map[string][]string, error) {
	tags, err := taglib.ReadTags(path)

	return tags, err
}

func (s *MetadataService) ReadCover(path string) ([]byte, error) {
	img, err := taglib.ReadImage(path)

	return img, err
}
