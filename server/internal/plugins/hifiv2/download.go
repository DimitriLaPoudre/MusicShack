package hifiv2

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/config"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
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

func (p *HifiV2) DownloadInfo(ctx context.Context, id string, quality string) (downloadItem, error) {
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

func downloadTidal(ctx context.Context, manifestRaw []byte, file *os.File) error {
	var manifest manifestTidal
	if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
		return fmt.Errorf("downloadTidal: unmarshal manisfestRaw: %w", err)
	}

	if len(manifest.Urls) <= 0 {
		return fmt.Errorf("downloadTidal: manifest.Urls[0]: %w", errors.New("not found"))
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", manifest.Urls[0], nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloadTidal: download url request: %w", err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			path := file.Name()
			file.Close()
			os.Remove(path)
			return fmt.Errorf("downloadTidal: read in file: %w", context.Canceled)
		default:
			n, err := resp.Body.Read(buf)
			if n > 0 {
				file.Write(buf[:n])
			}
			if err == io.EOF {
				return nil
			}
			if err != nil {
				path := file.Name()
				file.Close()
				os.Remove(path)
				return fmt.Errorf("downloadTidal: read in file: %w", err)
			}
		}
	}
}

func downloadMPD(ctx context.Context, manifest []byte, file *os.File) error {
	return nil
}

func (p *HifiV2) Download(ctx context.Context, userId uint, id string, quality string, status chan<- models.Status, data chan<- models.SongData) error {
	if quality == "" {
		quality = "LOSSLESS"
	}
	status <- models.StatusPending

	song, err := p.Song(ctx, id)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			status <- models.StatusCancel
		} else {
			status <- models.StatusFailed
		}
		return fmt.Errorf("HifiV2.Download: %w", err)
	}
	data <- song

	info, err := p.DownloadInfo(ctx, id, quality)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			status <- models.StatusCancel
		} else {
			status <- models.StatusFailed
		}
		return fmt.Errorf("HifiV2.Download: %w", err)
	}

	status <- models.StatusRunning

	filename := fmt.Sprintf("%s/%s/%s/%d - %s.", config.DOWNLOAD_FOLDER, song.Artist.Name, song.Album.Title, song.TrackNumber, song.Title)
	if quality == "HI_RES_LOSSLESS" || quality == "LOSSLESS" {
		filename += "flac"
	} else {
		filename += "mp4"
	}
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		status <- models.StatusFailed
		return fmt.Errorf("HifiV2.Download: os.MkdirAll: %w", err)

	}
	file, err := os.Create(filename)
	if err != nil {
		status <- models.StatusFailed
		return fmt.Errorf("HifiV2.Download: os.Create: %w", err)
	}
	defer file.Close()

	manifest, err := base64.StdEncoding.DecodeString(info.Manifest)
	if err != nil {
		return fmt.Errorf("HifiV2.Download: manifest decoding: %w", err)
	}

	switch info.ManifestMimeType {
	case "application/vnd.tidal.bts":
		err = downloadTidal(ctx, manifest, file)
	case "application/dash+xml":
		err = downloadMPD(ctx, manifest, file)
	default:
		err = errors.New("manifest type unknown")
	}

	if err != nil {
		if errors.Is(err, context.Canceled) {
			status <- models.StatusCancel
		} else {
			status <- models.StatusFailed
		}
		return fmt.Errorf("HifiV2.Download: %w", err)
	}

	if err := utils.FormatMetadata(filename, song); err != nil {
		status <- models.StatusFailed
		return fmt.Errorf("HifiV2.Download: %w", err)
	}

	status <- models.StatusDone

	return nil
}
