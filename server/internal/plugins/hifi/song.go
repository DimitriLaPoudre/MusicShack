package hifi

import (
	"context"
	"encoding/json"
	"maps"
	"net/http"
	"strings"

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

func (p *Hifi) Song(ctx context.Context, id string) (models.SongData, error) {
	apiInstance, err := repository.GetApiInstanceByApi(p.Name())
	if err != nil {
		return models.SongData{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiInstance.Url+"/track/?id="+id, nil)
	if err != nil {
		return models.SongData{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return models.SongData{}, err
	}
	defer resp.Body.Close()

	var items []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return models.SongData{}, err
	}

	data := make(map[string]any)
	for _, item := range items {
		maps.Copy(data, item)
	}

	var songData songData
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &songData,
		TagName: "mapstructure",
	})
	if err != nil {
		return models.SongData{}, err
	}
	if err := decoder.Decode(data); err != nil {
		return models.SongData{}, err
	}

	if songData.Album.CoverUrl != "" {
		songData.Album.CoverUrl = "https://resources.tidal.com/images/" + strings.ReplaceAll(songData.Album.CoverUrl, "-", "/") + "/640x640.jpg"
	}

	var normalizeSongData models.SongData
	decoder, err = mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           &normalizeSongData,
		TagName:          "useless",
		WeaklyTypedInput: true,
	})
	if err != nil {
		return models.SongData{}, err
	}
	if err := decoder.Decode(songData); err != nil {
		return models.SongData{}, err
	}

	return normalizeSongData, nil
}
