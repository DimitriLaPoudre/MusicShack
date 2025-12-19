package hifiv2

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
)

func getDownloadInfo(ctx context.Context, wg *sync.WaitGroup, apiUrl string, id string, quality string, ch chan<- downloadData) {
	defer wg.Done()

	url := apiUrl + "/track/?id=" + id
	if quality != "" {
		url += "&quality=" + quality
	} else {
		url += "&quality=LOSSLESS"
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data downloadData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}

	ch <- data
}

func (p *HifiV2) downloadInfo(ctx context.Context, id string, quality string) (downloadItem, error) {
	instances, err := repository.ListApiInstancesByApi(p.Name())
	if err != nil {
		return downloadItem{}, fmt.Errorf("HifiV2.DownloadInfo: %w", err)
	}

	routineCtx, routineCancel := context.WithCancel(context.Background())
	ch := make(chan downloadData)
	var wg sync.WaitGroup
	wg.Add(len(instances))
	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, instance := range instances {
		go getDownloadInfo(routineCtx, &wg, instance.Url, id, quality, ch)
	}

	var data downloadData
	select {
	case find, ok := <-ch:
		routineCancel()
		if !ok {
			return downloadItem{}, fmt.Errorf("HifiV2.DownloadInfo: %w", errors.New("can't be fetch"))
		}
		data = find
	case <-ctx.Done():
		routineCancel()
		return downloadItem{}, fmt.Errorf("HifiV2.DownloadInfo: %w", context.Canceled)
	}

	return data.Data, nil
}

func downloadTidal(ctx context.Context, manifestRaw []byte) (io.ReadCloser, error) {
	var manifest manifestTidal
	if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
		return nil, fmt.Errorf("downloadTidal: unmarshal manisfestRaw: %w", err)
	}

	if len(manifest.Urls) <= 0 {
		return nil, fmt.Errorf("downloadTidal: manifest.Urls[0]: %w", errors.New("not found"))
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", manifest.Urls[0], nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloadTidal: download url request: %w", err)
	}

	return resp.Body, nil
}

func downloadMPD(ctx context.Context, manifest []byte) (io.ReadCloser, error) {
	return nil, nil
}

func (p *HifiV2) Download(ctx context.Context, id string, quality string, data chan<- models.SongData) (io.ReadCloser, string, error) {
	if quality == "" {
		quality = "LOSSLESS"
	}

	song, err := p.Song(ctx, id)
	if err != nil {
		return nil, "", fmt.Errorf("HifiV2.Download: %w", err)
	}
	data <- song

	info, err := p.downloadInfo(ctx, id, quality)
	if err != nil {
		return nil, "", fmt.Errorf("HifiV2.Download: %w", err)
	}

	manifest, err := base64.StdEncoding.DecodeString(info.Manifest)
	if err != nil {
		return nil, "", fmt.Errorf("HifiV2.Download: manifest decoding: %w", err)
	}

	var reader io.ReadCloser
	switch info.ManifestMimeType {
	case "application/vnd.tidal.bts":
		reader, err = downloadTidal(ctx, manifest)
	case "application/dash+xml":
		reader, err = downloadMPD(ctx, manifest)
	default:
		err = errors.New("manifest type unknown")
	}

	if err != nil {
		return nil, "", fmt.Errorf("HifiV2.Download: %w", err)
	}

	var extension string
	if quality == "HI_RES_LOSSLESS" || quality == "LOSSLESS" {
		extension = "flac"
	} else {
		extension = "mp4"
	}

	return reader, extension, nil
}
