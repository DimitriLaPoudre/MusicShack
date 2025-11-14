package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"maps"
	"net/http"
	"strings"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/mitchellh/mapstructure"
)

type songData struct {
	Id           uint   `mapstructure:"id"`
	Title        string `mapstructure:"title"`
	Duration     uint   `mapstructure:"duration"`
	ReleaseDate  string `mapstructure:"streamStartDate"`
	TrackNumber  uint   `mapstructure:"trackNumber"`
	VolumeNumber uint   `mapstructure:"volumeNumber"`
	AudioQuality string `mapstructure:"audioQuality"`
	Artist       struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artist"`
	Artists []struct {
		Id   uint   `mapstructure:"id"`
		Name string `mapstructure:"name"`
	} `mapstructure:"artists"`
	Album struct {
		Id       uint   `mapstructure:"id"`
		Title    string `mapstructure:"title"`
		CoverUrl string `mapstructure:"cover"`
	} `mapstructure:"album"`
	BitDepth    uint   `mapstructure:"bitDepth"`
	SampleRate  uint   `mapstructure:"sampleRate"`
	DownloadUrl string `mapstructure:"OriginalTrackUrl"`
}

func (p *Hifi) getSong(ctx context.Context, wg *sync.WaitGroup, urlApi string, ch chan<- songData, id string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlApi+"/track/?id="+id, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return
	}

	tmp := make(map[string]any)
	for _, item := range items {
		maps.Copy(tmp, item)
	}

	var data songData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &data,
		TagName: "mapstructure",
	})
	if err != nil {
		return
	}
	if err := decoder.Decode(tmp); err != nil {
		return
	}

	ch <- data
}

func (p *Hifi) Song(ctx context.Context, id string) (models.SongData, error) {
	apiInstances, err := repository.ListApiInstancesByApi(p.Name())
	if err != nil {
		return models.SongData{}, err
	}

	var data songData
	routineCtx, cancel := context.WithCancel(context.Background())
	ch := make(chan songData)
	var wg sync.WaitGroup

	wg.Add(len(apiInstances))
	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, instance := range apiInstances {
		go p.getSong(routineCtx, &wg, instance.Url, ch, id)
	}
	select {
	case find, ok := <-ch:
		if !ok {
			return models.SongData{}, errors.New("Song not found")
		}
		cancel()
		data = find
	case <-ctx.Done():
		cancel()
	}

	if data.Album.CoverUrl != "" {
		data.Album.CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(data.Album.CoverUrl, "-", "/") + "/640x640.jpg"
	}

	var normalizeSongData models.SongData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeSongData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.SongData{}, err
	}
	if err := decoder.Decode(data); err != nil {
		return models.SongData{}, err
	}

	return normalizeSongData, nil
}
