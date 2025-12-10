package hifiv2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
	"github.com/mitchellh/mapstructure"
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

	var data songData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) Song(ctx context.Context, id string) (models.SongData, error) {
	apiInstances, err := repository.ListApiInstancesByApi(p.Name())
	if err != nil {
		return models.SongData{}, fmt.Errorf("HifiV2.Song: %w", err)
	}

	routineCtx, routineCancel := context.WithCancel(context.Background())
	ch := make(chan songData)
	var wg sync.WaitGroup
	wg.Add(len(apiInstances))
	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, instance := range apiInstances {
		go getSong(routineCtx, &wg, instance.Url, ch, id)
	}

	var data songData
	select {
	case find, ok := <-ch:
		routineCancel()
		if !ok {
			return models.SongData{}, fmt.Errorf("HifiV2.Song: %w", errors.New("can't be fetch"))
		}
		data = find
	case <-ctx.Done():
		routineCancel()
	}

	data.Data.Album.CoverUrl = utils.GetImageURL(data.Data.Album.CoverUrl, 640)

	var normalizeSongData models.SongData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeSongData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.SongData{}, fmt.Errorf("HifiV2.Song: %w", err)
	}
	if err := decoder.Decode(data.Data); err != nil {
		return models.SongData{}, fmt.Errorf("HifiV2.Song: %w", err)
	}

	fmt.Printf("%#v\n", normalizeSongData)

	return normalizeSongData, nil
}
