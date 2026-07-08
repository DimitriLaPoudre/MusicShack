package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	lib_url "net/url"
	"strconv"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	hifi_utils "github.com/Ascension-EIP/Ascension/apps/server/internal/outbound/plugin/hifi/utils"
	"golang.org/x/time/rate"
)

func fetchAlbum(ctx context.Context, limiter *rate.Limiter, url string, id string) (albumData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := hifi_utils.Fetch(ctx, url+"/album/?id="+lib_url.QueryEscape(id), limiter)
	if err != nil {
		return albumData{}, fmt.Errorf("fetchAlbum: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return albumData{}, fmt.Errorf("fetchAlbum: http: %w", errors.New(resp.Status))
	}

	var data albumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return albumData{}, fmt.Errorf("fetchAlbum: json.Decode: %w", err)
	}

	return data, nil
}

func getAlbum(ctx context.Context, limiter *rate.Limiter, url string, id string) (albumData, error) {
	album, err := fetchAlbum(ctx, limiter, url, id)
	if err != nil {
		return albumData{}, nil
	}

	return album, nil
}

func (p *Hifi) Album(ctx context.Context, url string, id string) (model.Album, error) {
	album, err := getAlbum(ctx, p.limiter, url, id)
	if err != nil {
		return model.Album{}, fmt.Errorf("Hifi.Album: %w", err)
	}

	releaseDate, err := time.Parse(StreamStartDateLayout, album.Data.ReleaseDate)
	if err != nil {
		p.l.Warn().Msg(fmt.Sprintf("Hifi.Album: time.Parse(%s): %s", album.Data.ReleaseDate, err.Error()))
	}

	audioQuality := LOW
	switch album.Data.AudioQuality {
	case "LOW":
		audioQuality = LOW
	case "HIGH":
		audioQuality = HIGH
	case "LOSSLESS":
		audioQuality = LOSSLESS
	default:
		audioQuality = LOW
	}
	for _, quality := range album.Data.MediaMetadata.Tags {
		switch quality {
		case "HIRES_LOSSLESS":
			audioQuality = HIRES
		case "LOSSLESS", "DOLBY_ATMOS":
			if audioQuality != HIRES {
				audioQuality = LOSSLESS
			}
		}
	}

	songs := []model.Song{}
	for _, item := range album.Data.Items {
		song := item.Item

		audioQuality := LOW
		switch song.AudioQuality {
		case "LOW":
			audioQuality = LOW
		case "HIGH":
			audioQuality = HIGH
		case "LOSSLESS":
			audioQuality = LOSSLESS
		default:
			audioQuality = LOW
		}
		for _, quality := range song.MediaMetadata.Tags {
			switch quality {
			case "HIRES_LOSSLESS":
				audioQuality = HIRES
			case "LOSSLESS", "DOLBY_ATMOS":
				if audioQuality != HIRES {
					audioQuality = LOSSLESS
				}
			}
		}

		artists := []model.Artist{}
		for _, artist := range song.Artists {
			artists = append(artists, model.Artist{
				Id:   strconv.FormatUint(uint64(artist.Id), 10),
				Name: artist.Name,
			})
		}

		songs = append(songs, model.Song{
			Id:           strconv.FormatUint(uint64(song.Id), 10),
			Title:        song.Title,
			Duration:     song.Duration,
			TrackNumber:  song.TrackNumber,
			VolumeNumber: song.VolumeNumber,
			AudioQuality: audioQuality,
			Explicit:     song.Explicit,
			Isrc:         song.Isrc,
			Artists:      artists,
		})
	}

	artists := []model.Artist{}
	for _, artist := range album.Data.Artists {
		artists = append(artists, model.Artist{
			Id:   strconv.FormatUint(uint64(artist.Id), 10),
			Name: artist.Name,
		})
	}

	return model.Album{
		Id:            strconv.FormatUint(uint64(album.Data.Id), 10),
		Title:         album.Data.Title,
		Duration:      album.Data.Duration,
		ReleaseDate:   releaseDate,
		NumberTracks:  album.Data.NumberOfTracks,
		NumberVolumes: album.Data.NumberOfVolumes,
		CoverUrl:      hifi_utils.GetImageURL(album.Data.CoverUrl, 640),
		AudioQuality:  audioQuality,
		Explicit:      album.Data.Explicit,
		Songs:         songs,
		Artists:       artists,
	}, nil
}
