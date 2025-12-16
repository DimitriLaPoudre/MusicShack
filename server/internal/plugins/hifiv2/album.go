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

func getAlbum(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- albumData, id string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/album/?id="+id, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data albumData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Album(ctx context.Context, id string) (models.AlbumData, error) {
	apiInstances, err := repository.ListApiInstancesByApi(p.Name())
	if err != nil {
		return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", err)
	}

	routineCtx, cancel := context.WithCancel(context.Background())
	ch := make(chan albumData)
	var wg sync.WaitGroup

	wg.Add(len(apiInstances))
	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, instance := range apiInstances {
		go getAlbum(routineCtx, &wg, instance.Url, ch, id)
	}

	var data albumData
	select {
	case find, ok := <-ch:
		cancel()
		if !ok {
			return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", errors.New("can't fetch"))
		} else {
			data = find
		}
	case <-ctx.Done():
		cancel()
		return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", context.Canceled)
	}

	var normalizeAlbumData models.AlbumData

	normalizeAlbumData.Limit = data.Data.Limit
	normalizeAlbumData.Offset = data.Data.Offset
	normalizeAlbumData.NumberSongs = data.Data.TotalNumberOfItems
	for _, wrappedSong := range data.Data.Items {
		dirtySong := wrappedSong.Item
		var artists []struct {
			Id   string
			Name string
		}
		for _, artist := range dirtySong.Artists {
			artists = append(artists, struct {
				Id   string
				Name string
			}{
				strconv.FormatUint(uint64(artist.Id), 10),
				artist.Name,
			})
		}
		song := struct {
			Id           string
			Title        string
			Duration     uint
			TrackNumber  uint
			VolumeNumber uint
			Artists      []struct {
				Id   string
				Name string
			}
		}{strconv.FormatUint(uint64(dirtySong.Id), 10),
			dirtySong.Title,
			dirtySong.Duration,
			dirtySong.TrackNumber,
			dirtySong.VolumeNumber,
			artists,
		}
		normalizeAlbumData.Duration += song.Duration
		normalizeAlbumData.Songs = append(normalizeAlbumData.Songs, song)
	}
	if len(normalizeAlbumData.Songs) == 0 {
		return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %v: %w", normalizeAlbumData, errors.New("songs not found"))
	}

	firstSong := data.Data.Items[0].Item

	normalizeAlbumData.Id = strconv.FormatUint(uint64(firstSong.Album.Id), 10)
	normalizeAlbumData.Title = firstSong.Album.Title
	normalizeAlbumData.ReleaseDate = firstSong.ReleaseDate[:10]
	normalizeAlbumData.CoverUrl = utils.GetImageURL(firstSong.Album.CoverUrl, 640)

	return normalizeAlbumData, nil
}
