package hifi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
	pkg_io "github.com/DimitriLaPoudre/MusicShack/internal/pkg/io"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/network"
)

func getDownloadInfo(ctx context.Context, limiters *sync.Map, urls []string, id string, quality string) (downloadResponse, error) {
	path := "/track/?id=" + url.QueryEscape(id)
	if quality != "" {
		path += "&quality=" + url.QueryEscape(quality)
	}

	downloadInfo, err := hifi_utils.FetchTypeSequential[downloadResponse](ctx, urls, path, limiters)
	if err != nil {
		return downloadResponse{}, fmt.Errorf("fetch download info with url list: %w", err)
	}

	return downloadInfo, nil
}

func downloadTidal(ctx context.Context, manifestRaw []byte) (io.ReadCloser, error) {
	var manifest manifestTidal
	if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
		return nil, fmt.Errorf("json unmarshal manifest: %w", err)
	}

	if len(manifest.Urls) <= 0 {
		return nil, fmt.Errorf("manifest first url: %w", errors.New("not found"))
	}

	resp, err := network.Fetch(ctx, manifest.Urls[0])
	if err != nil {
		return nil, fmt.Errorf("fetch download: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch download response status: %w", errors.New(resp.Status))
	}

	return resp.Body, nil
}

func downloadMPD(ctx context.Context, manifest []byte) (io.ReadCloser, error) {
	var mpd manifestMPD
	if err := xml.Unmarshal(manifest, &mpd); err != nil {
		return nil, fmt.Errorf("xml unmarshal manifest: %w", err)
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
		resp, err := network.Fetch(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("fetch download segment: %w", err)
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			resp.Body.Close()
			return nil, fmt.Errorf("fetch download segment response status: %w", errors.New(resp.Status))
		}

		readers = append(readers, resp.Body)
	}

	fullReader := pkg_io.MultiReadCloser(readers...)

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
		reader.Close()
		return nil, fmt.Errorf("open stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		reader.Close()
		return nil, fmt.Errorf("open stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		reader.Close()
		return nil, fmt.Errorf("start ffmpeg remux command: %w", err)
	}

	go func() {
		if _, err := io.Copy(stdin, reader); err != nil {
			slog.Error("copy from source stream to stdin: %w", err)
		}
		stdin.Close()
		reader.Close()
	}()

	newReader, writer := io.Pipe()

	go func() {
		if _, err := io.Copy(writer, stdout); err != nil {
			writer.CloseWithError(err)
		}

		if err := cmd.Wait(); err != nil {
			writer.CloseWithError(err)
		} else {
			writer.Close()
		}
	}()
	return newReader, nil
}

func (p *Hifi) Download(ctx context.Context, instances []model.Instance, id string, hiRes bool) (io.ReadCloser, string, error) {
	quality := "LOSSLESS"
	if hiRes {
		quality = "HI_RES_LOSSLESS"
	}

	urls := hifi_utils.InstancesToUrls(instances)

	downloadInfo, err := getDownloadInfo(ctx, &p.limiters, urls, id, quality)
	if err != nil {
		return nil, "", err
	}

	manifest, err := base64.StdEncoding.DecodeString(downloadInfo.Data.Manifest)
	if err != nil {
		return nil, "", fmt.Errorf("decoding song %s download manifest: %w", id, err)
	}

	var extension string
	switch downloadInfo.Data.AudioQuality {
	case AudioQualityHIRES:
		if quality != "HI_RES_LOSSLESS" {
			return nil, "", fmt.Errorf("download song %s: %w", id, errors.New("audio quality received too high"))
		}
		extension = "flac"
	case AudioQualityLOSSLESS:
		extension = "flac"
	case AudioQualityHIGH:
		extension = "m4a"
	}

	var reader io.ReadCloser
	switch downloadInfo.Data.ManifestMimeType {
	case "application/vnd.tidal.bts":
		reader, err = downloadTidal(ctx, manifest)
	case "application/dash+xml":
		reader, err = downloadMPD(ctx, manifest)
	default:
		err = errors.New("manifest type unknown")
	}
	if err != nil {
		return nil, "", fmt.Errorf("downloading song %s: %w", id, err)
	}

	switch downloadInfo.Data.AudioQuality {
	case AudioQualityHIRES, AudioQualityLOSSLESS:
		reader, err = remuxM4AtoFLAC(reader)
		if err != nil {
			return nil, "", fmt.Errorf("remux from M4A to FLAC: %w", err)
		}
	}

	return reader, extension, nil
}
