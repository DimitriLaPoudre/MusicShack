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
			return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", errors.New("can't be fetch"))
		} else {
			data = find
		}
	case <-ctx.Done():
		cancel()
		return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", errors.New("context canceled"))
	}

	var normalizeAlbumData models.AlbumData

	normalizeAlbumData.Limit = data.Data.Limit
	normalizeAlbumData.Offset = data.Data.Offset
	normalizeAlbumData.NumberSongs = data.Data.TotalNumberOfItems
	for _, wrappedSong := range data.Data.Items {
		dirtySong := wrappedSong.Item
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
			[]struct {
				Id   string
				Name string
			}{},
		}
		normalizeAlbumData.Duration += song.Duration
		normalizeAlbumData.Songs = append(normalizeAlbumData.Songs, song)
	}
	if len(normalizeAlbumData.Songs) == 0 {
		return models.AlbumData{}, fmt.Errorf("HifiV2.Album: %w", errors.New("can't be formatted"))
	}

	firstSong := normalizeAlbumData.Songs[0]

	normalizeAlbumData.Id = firstSong.Id
	normalizeAlbumData.Title = firstSong.Title
	normalizeAlbumData.ReleaseDate = data.Data.Items[0].Item.ReleaseDate
	normalizeAlbumData.CoverUrl = utils.GetImageURL(data.Data.Items[0].Item.Album.CoverUrl, 640)

	return normalizeAlbumData, nil
}
