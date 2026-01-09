package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func fetchSearchSong(ctx context.Context, urlApi string, song string) (searchSongData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?s="+url.QueryEscape(song), nil)
	if err != nil {
		return searchSongData{}, fmt.Errorf("fetchSearchSong: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return searchSongData{}, fmt.Errorf("fetchSearchSong: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return searchSongData{}, fmt.Errorf("fetchSearchSong: %w", errors.New("http error "+strconv.FormatInt(int64(resp.StatusCode), 10)))
	}

	var data searchSongData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return searchSongData{}, fmt.Errorf("fetchSearchSong: json.Decode: %w", err)
	}

	return data, nil
}

func getSearchSong(ctx context.Context, instances []models.Instance, song string) (searchSongData, error) {
	type res struct {
		data searchSongData
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchSearchSong(ctx, url, song)
			ch <- res{data: data, err: err}
		}(instance.Url)
	}

	var lastErr error
	for range instances {
		select {
		case res := <-ch:
			if res.err == nil {
				return res.data, nil
			}
			lastErr = res.err
		case <-ctx.Done():
			return searchSongData{}, ctx.Err()
		}
	}
	return searchSongData{}, fmt.Errorf("getSearchSong: %w", lastErr)
}

func fetchSearchAlbum(ctx context.Context, urlApi string, album string) (searchAlbumData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?al="+url.QueryEscape(album), nil)
	if err != nil {
		return searchAlbumData{}, fmt.Errorf("fetchSearchAlbum: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return searchAlbumData{}, fmt.Errorf("fetchSearchAlbum: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return searchAlbumData{}, fmt.Errorf("fetchSearchAlbum: %w", errors.New("http error "+strconv.FormatInt(int64(resp.StatusCode), 10)))
	}

	var data searchAlbumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return searchAlbumData{}, fmt.Errorf("fetchSearchAlbum: json.Decode: %w", err)
	}

	return data, nil
}

func getSearchAlbum(ctx context.Context, instances []models.Instance, album string) (searchAlbumData, error) {
	type res struct {
		data searchAlbumData
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchSearchAlbum(ctx, url, album)
			ch <- res{data: data, err: err}
		}(instance.Url)
	}

	var lastErr error
	for range instances {
		select {
		case res := <-ch:
			if res.err == nil {
				return res.data, nil
			}
			lastErr = res.err
		case <-ctx.Done():
			return searchAlbumData{}, ctx.Err()
		}
	}
	return searchAlbumData{}, fmt.Errorf("getSearchAlbum: %w", lastErr)
}

func fetchSearchArtist(ctx context.Context, urlApi string, artist string) (searchArtistData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/search/?a="+url.QueryEscape(artist), nil)
	if err != nil {
		return searchArtistData{}, fmt.Errorf("fetchSearchArtist: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return searchArtistData{}, fmt.Errorf("fetchSearchArtist: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return searchArtistData{}, fmt.Errorf("fetchSearchArtist: %w", errors.New("http error "+strconv.FormatInt(int64(resp.StatusCode), 10)))
	}

	var data searchArtistData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return searchArtistData{}, fmt.Errorf("fetchSearchArtist: json.Decode: %w", err)
	}

	return data, nil
}

func getSearchArtist(ctx context.Context, instances []models.Instance, artist string) (searchArtistData, error) {
	type res struct {
		data searchArtistData
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchSearchArtist(ctx, url, artist)
			ch <- res{data: data, err: err}
		}(instance.Url)
	}

	var lastErr error
	for range instances {
		select {
		case res := <-ch:
			if res.err == nil {
				return res.data, nil
			}
			lastErr = res.err
		case <-ctx.Done():
			return searchArtistData{}, ctx.Err()
		}
	}
	return searchArtistData{}, fmt.Errorf("getSearchArtist: %w", lastErr)
}

func getSearchData(ctx context.Context, instances []models.Instance, song, album, artist string) (searchSongData, searchAlbumData, searchArtistData, error) {
	type res struct {
		data any
		err  error
	}

	ch := make(chan res, 3)
	go func() {
		data, err := getSearchSong(ctx, instances, song)
		ch <- res{data: data, err: err}
	}()
	go func() {
		data, err := getSearchAlbum(ctx, instances, album)
		ch <- res{data: data, err: err}
	}()
	go func() {
		data, err := getSearchArtist(ctx, instances, artist)
		ch <- res{data: data, err: err}
	}()

	var songData searchSongData
	var albumData searchAlbumData
	var artistData searchArtistData
	for range 3 {
		res := <-ch
		switch v := res.data.(type) {
		case searchSongData:
			songData = v
		case searchAlbumData:
			albumData = v
		case searchArtistData:
			artistData = v
		}
	}

	return songData, albumData, artistData, nil
}

func (p *Hifi) Search(ctx context.Context, userId uint, song, album, artist string) (models.SearchData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.SearchData{}, fmt.Errorf("Hifi.Search: %w", err)
	}
	if len(instances) == 0 {
		return models.SearchData{}, fmt.Errorf("Hifi.Search: %w", errors.New("not found"))
	}

	songData, albumData, artistData, err := getSearchData(ctx, instances, song, album, artist)
	if err != nil {
		return models.SearchData{}, fmt.Errorf("Hifi.Search: %w", err)
	}

	result := models.SearchData{
		Songs:   make([]models.SearchDataSong, 0),
		Albums:  make([]models.SearchDataAlbum, 0),
		Artists: make([]models.SearchDataArtist, 0),
	}

	if len(songData.Data.Songs) != 0 {
		for _, rawSong := range songData.Data.Songs {
			song := models.SearchDataSong{
				Id:           strconv.FormatUint(uint64(rawSong.Id), 10),
				Title:        rawSong.Title,
				Duration:     rawSong.Duration,
				AudioQuality: models.QualityHigh,
				Popularity:   rawSong.Popularity,
				Explicit:     rawSong.Explicit,
				Artists:      make([]models.SongDataArtist, 0),
				Album: models.SongDataAlbum{
					Id:       strconv.FormatUint(uint64(rawSong.Album.Id), 10),
					Title:    rawSong.Album.Title,
					CoverUrl: utils.GetImageURL(rawSong.Album.CoverUrl, 1280),
				},
			}

			for _, quality := range rawSong.MediaMetadata.Tags {
				switch quality {
				case "HIRES_LOSSLESS":
					song.AudioQuality = max(song.AudioQuality, models.QualityHiresLossless)
				case "LOSSLESS", "DOLBY_ATMOS":
					song.AudioQuality = max(song.AudioQuality, models.QualityLossless)
				}
			}

			for _, artist := range rawSong.Artists {
				song.Artists = append(song.Artists, models.SongDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			result.Songs = append(result.Songs, song)
		}
	}

	if len(albumData.Data.Albums.Albums) != 0 {
		for _, rawAlbum := range albumData.Data.Albums.Albums {
			album := models.SearchDataAlbum{
				Id:           strconv.FormatUint(uint64(rawAlbum.Id), 10),
				Title:        rawAlbum.Title,
				Duration:     rawAlbum.Duration,
				CoverUrl:     utils.GetImageURL(rawAlbum.CoverUrl, 1280),
				AudioQuality: models.QualityHigh,
				Explicit:     rawAlbum.Explicit,
				Popularity:   rawAlbum.Popularity,
				Artists:      make([]models.AlbumDataArtist, 0),
			}

			for _, quality := range rawAlbum.MediaMetadata.Tags {
				switch quality {
				case "HIRES_LOSSLESS":
					album.AudioQuality = max(album.AudioQuality, models.QualityHiresLossless)
				case "LOSSLESS", "DOLBY_ATMOS":
					album.AudioQuality = max(album.AudioQuality, models.QualityLossless)
				}
			}

			for _, artist := range rawAlbum.Artists {
				album.Artists = append(album.Artists, models.AlbumDataArtist{
					Id:   strconv.FormatUint(uint64(artist.Id), 10),
					Name: artist.Name,
				})
			}

			result.Albums = append(result.Albums, album)
		}
	}

	if len(artistData.Data.Artists.Artists) != 0 {
		for _, rawArtist := range artistData.Data.Artists.Artists {
			result.Artists = append(result.Artists, models.SearchDataArtist{
				Id:         strconv.FormatUint(uint64(rawArtist.Id), 10),
				Name:       rawArtist.Name,
				PictureUrl: utils.GetImageURL(rawArtist.PictureUrl, 750),
				Popularity: rawArtist.Popularity,
			})
		}
	}

	return result, nil
}
