package hifi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func fetchDownloadInfo(ctx context.Context, url2 string, id string, quality string) (downloadData, error) {
	url2 += "/track/?id=" + url.QueryEscape(id)
	if quality != "" {
		url2 += "&quality=" + url.QueryEscape(quality)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := utils.Fetch(ctx, url2)
	if err != nil {
		return downloadData{}, fmt.Errorf("fetchDownloadInfo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return downloadData{}, fmt.Errorf("fetchDownloadInfo: http: %w", errors.New(resp.Status))
	}

	var data downloadData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return downloadData{}, fmt.Errorf("fetchDownloadInfo: json.Decode: %w", err)
	}

	return data, nil
}

func getDownloadInfo(ctx context.Context, instances []models.Instance, id string, quality string) (downloadData, error) {
	type res struct {
		data downloadData
		err  error
	}

	ch := make(chan res, len(instances))
	for _, instance := range instances {
		go func(url string) {
			data, err := fetchDownloadInfo(ctx, url, id, quality)
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
			return downloadData{}, ctx.Err()
		}
	}
	return downloadData{}, fmt.Errorf("getDownloadInfo: %w", lastErr)
}

func downloadTidal(ctx context.Context, manifestRaw []byte) (io.ReadCloser, error) {
	var manifest manifestTidal
	if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
		return nil, fmt.Errorf("downloadTidal: json.Unmarshal: %w", err)
	}

	if len(manifest.Urls) <= 0 {
		return nil, fmt.Errorf("downloadTidal: manifest.Urls[0]: %w", errors.New("not found"))
	}

	resp, err := utils.Fetch(ctx, manifest.Urls[0])
	if err != nil {
		return nil, fmt.Errorf("downloadTidal: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("downloadTidal: http: %w", errors.New(resp.Status))
	}

	return resp.Body, nil
}

func downloadMPD(ctx context.Context, manifest []byte) (io.ReadCloser, error) {
	var mpd manifestMPD
	if err := xml.Unmarshal(manifest, &mpd); err != nil {
		return nil, fmt.Errorf("downloadMPD: xml.Unmarshal: %w", err)
	}

	rep := mpd.Periods[0].
		AdaptationSets[0].
		Representations[0]

	tmpl := rep.SegmentTemplate

	segments := []string{tmpl.Initialization}

	n := tmpl.StartNumber
	for _, s := range tmpl.Timeline.Segments {
		repeat := max(s.R, 0)
		for i := 0; i <= repeat; i++ {
			url := strings.Replace(
				tmpl.Media,
				"$Number$",
				strconv.Itoa(n),
				1,
			)
			segments = append(segments, url)
			n++
		}
	}

	var readers []io.ReadCloser

	for _, url := range segments {
		resp, err := utils.Fetch(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("fetchAlbum: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			defer resp.Body.Close()
			return nil, fmt.Errorf("fetchAlbum: http: %w", errors.New(resp.Status))
		}

		readers = append(readers, resp.Body)
	}

	fullReader := utils.MultiReadCloser(readers...)

	return fullReader, nil
}

func remuxM4AtoFLAC(reader io.ReadCloser) (io.ReadCloser, error) {
	cmd := exec.Command("ffmpeg",
		"-nostdin",
		"-fflags", "+genpts",
		"-i", "pipe:0",
		"-map", "0:a:0",
		"-map_metadata", "0",
		"-c:a", "copy",
		"-f", "flac",
		"pipe:1")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		_ = reader.Close()
		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = reader.Close()
		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
	}

	if err := cmd.Start(); err != nil {
		_ = reader.Close()
		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
	}

	go func() {
		if _, err := io.Copy(stdin, reader); err != nil {
			log.Println("io.Copy(stdin, reader): %w", err)
		}
		_ = stdin.Close()
		_ = reader.Close()
	}()

	newReader, writer := io.Pipe()

	go func() {
		if _, err := io.Copy(writer, stdout); err != nil {
			writer.CloseWithError(err)
		}

		if err := cmd.Wait(); err != nil {
			_ = writer.CloseWithError(err)
		} else {
			_ = writer.Close()
		}
	}()
	return newReader, nil
}

func (p *Hifi) Download(ctx context.Context, userId uint, id string) (io.ReadCloser, string, error) {
	quality := "HI_RES_LOSSLESS"
	if user, err := repository.GetUserByID(userId); err != nil {
		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
	} else if !user.HiRes {
		quality = "LOSSLESS"
	}

	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
	if err != nil {
		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
	}

	info, err := getDownloadInfo(ctx, instances, id, quality)
	if err != nil {
		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
	}

	manifest, err := base64.StdEncoding.DecodeString(info.Data.Manifest)
	if err != nil {
		return nil, "", fmt.Errorf("Hifi.Download: base64.StdEncoding.DecodeString: %w", err)
	}

	var reader io.ReadCloser
	switch info.Data.ManifestMimeType {
	case "application/vnd.tidal.bts":
		reader, err = downloadTidal(ctx, manifest)
	case "application/dash+xml":
		reader, err = downloadMPD(ctx, manifest)
	default:
		err = errors.New("manifest type unknown")
	}
	if err != nil {
		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
	}

	var extension string
	switch info.Data.AudioQuality {
	case "HI_RES_LOSSLESS":
		if quality != "HI_RES_LOSSLESS" {
			return nil, "", fmt.Errorf("Hifi.Download: %w", errors.New("audio quality received not conform"))
		}
		reader, err = remuxM4AtoFLAC(reader)
		if err != nil {
			return nil, "", fmt.Errorf("Hifi.Download: %w", err)
		}
		extension = "flac"
	case "LOSSLESS":
		extension = "flac"
	case "HIGH":
		extension = "m4a"
	}

	return reader, extension, nil
}
