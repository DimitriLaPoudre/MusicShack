package hifiv2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func getSong(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- songData, id string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/info/?id="+id, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return
	}

	var data songData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Song(ctx context.Context, userId uint, id string) (models.SongData, error) {
	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return models.SongData{}, fmt.Errorf("HifiV2.Song: %w", err)
	}

	routineCtx, routineCancel := context.WithCancel(context.Background())
	ch := make(chan songData)
	var wg sync.WaitGroup
	wg.Add(len(instances))
	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, instance := range instances {
		go getSong(routineCtx, &wg, instance.Url, ch, id)
	}

	var data songData
	select {
	case find, ok := <-ch:
		routineCancel()
		if !ok {
			return models.SongData{}, fmt.Errorf("HifiV2.Song: %w", errors.New("can't fetch"))
		}
		data = find
	case <-ctx.Done():
		routineCancel()
		return models.SongData{}, fmt.Errorf("HifiV2.Song: %w", context.Canceled)
	}

	var normalizeSongData models.SongData
	{

		normalizeSongData.Id = strconv.FormatUint(uint64(data.Data.Id), 10)
		normalizeSongData.Title = data.Data.Title
		normalizeSongData.Duration = data.Data.Duration
		normalizeSongData.ReleaseDate = data.Data.ReleaseDate[:10]
		normalizeSongData.TrackNumber = data.Data.TrackNumber
		normalizeSongData.VolumeNumber = data.Data.VolumeNumber
		switch data.Data.AudioQuality {
		case "HI_RES_LOSSLESS":
			normalizeSongData.AudioQuality = 4
		case "LOSSLESS":
			normalizeSongData.AudioQuality = 3
		case "HIGH":
			normalizeSongData.AudioQuality = 2
		case "LOW":
			normalizeSongData.AudioQuality = 1
		}

		normalizeSongData.Popularity = data.Data.Popularity
		normalizeSongData.Isrc = data.Data.Isrc

		artists := make([]models.SongDataArtist, 0)
		for _, artist := range data.Data.Artists {
			artists = append(artists, models.SongDataArtist{
				Id:   strconv.FormatUint(uint64(artist.Id), 10),
				Name: artist.Name,
			})
		}
		normalizeSongData.Artists = artists

		album := models.SongDataAlbum{
			Id:       strconv.FormatUint(uint64(data.Data.Album.Id), 10),
			Title:    data.Data.Album.Title,
			CoverUrl: utils.GetImageURL(data.Data.Album.CoverUrl, 640),
		}
		normalizeSongData.Album = album
	}

	return normalizeSongData, nil
}
