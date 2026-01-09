package hifi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
	"github.com/DimitriLaPoudre/MusicShack/server/internal/repository"
)

func fetchDownloadInfo(ctx context.Context, url string, id string, quality string) (downloadData, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	url += "/track/?id=" + id
	if quality != "" {
		url += "&quality=" + quality
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return downloadData{}, fmt.Errorf("fetchDownloadInfo: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return downloadData{}, fmt.Errorf("fetchDownloadInfo: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return downloadData{}, fmt.Errorf("fetchDownloadInfo: %w", errors.New("http error "+strconv.FormatInt(int64(resp.StatusCode), 10)))
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
	req, err := http.NewRequestWithContext(ctx, "GET", manifest.Urls[0], nil)
	if err != nil {
		return nil, fmt.Errorf("downloadTidal: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("downloadTidal: http.DefaultClient.Do: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("downloadTidal: %w", errors.New("http error "+strconv.FormatInt(int64(resp.StatusCode), 10)))
	}

	return resp.Body, nil
}

type concatReadCloser struct {
	readers []io.ReadCloser
	current int
}

func NewConcatReadCloser(readers ...io.ReadCloser) io.ReadCloser {
	return &concatReadCloser{readers: readers}
}

func (c *concatReadCloser) Read(p []byte) (int, error) {
	for c.current < len(c.readers) {
		n, err := c.readers[c.current].Read(p)
		if err == io.EOF {
			c.readers[c.current].Close()
			c.current++
			continue
		}
		return n, err
	}
	return 0, io.EOF
}

func (c *concatReadCloser) Close() error {
	var firstErr error
	for i := c.current; i < len(c.readers); i++ {
		if err := c.readers[i].Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
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
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("downloadMPD: http.NewRequestWithContext: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("downloadMPD: http.DefaultClient.Do: %w", err)
		}

		if resp.StatusCode >= 400 {
			defer resp.Body.Close()
			return nil, fmt.Errorf("downloadMPD: %w", errors.New("http error "+strconv.FormatInt(int64(resp.StatusCode), 10)))
		}

		readers = append(readers, resp.Body)
	}

	fullReader := NewConcatReadCloser(readers...)

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
		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
	}

	go func() {
		defer stdin.Close()
		io.Copy(stdin, reader)
	}()

	newReader, writer := io.Pipe()

	go func() {
		io.Copy(writer, stdout)

		if err := cmd.Wait(); err != nil {
			writer.CloseWithError(err)
		} else {
			writer.Close()
		}
	}()
	return newReader, nil
}

func (p *Hifi) Download(ctx context.Context, userId uint, id string, quality string) (io.ReadCloser, string, error) {
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
		return nil, "", fmt.Errorf("Hifi.Download: manifest decoding: %w", err)
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
		reader, err = remuxM4AtoFLAC(reader)
		if err != nil {
			return nil, "", fmt.Errorf("Hifi.Download: %w", err)
		}
		extension = "flac"
	case "LOSSLESS":
		extension = "flac"
	case "HIGH":
		extension = "m4a"
	case "LOW":
		extension = "m4a"
	}

	return reader, extension, nil
}
