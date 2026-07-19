package hifi

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
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

// func downloadTidal(ctx context.Context, manifestRaw []byte) (io.ReadCloser, error) {
// 	var manifest manifestTidal
// 	if err := json.Unmarshal(manifestRaw, &manifest); err != nil {
// 		return nil, fmt.Errorf("downloadTidal: json.Unmarshal: %w", err)
// 	}
//
// 	if len(manifest.Urls) <= 0 {
// 		return nil, fmt.Errorf("downloadTidal: manifest.Urls[0]: %w", errors.New("not found"))
// 	}
//
// 	resp, err := utils.Fetch(ctx, manifest.Urls[0])
// 	if err != nil {
// 		return nil, fmt.Errorf("downloadTidal: %w", err)
// 	}
//
// 	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 		return nil, fmt.Errorf("downloadTidal: http: %w", errors.New(resp.Status))
// 	}
//
// 	return resp.Body, nil
// }
//
// func downloadMPD(ctx context.Context, manifest []byte) (io.ReadCloser, error) {
// 	var mpd manifestMPD
// 	if err := xml.Unmarshal(manifest, &mpd); err != nil {
// 		return nil, fmt.Errorf("downloadMPD: xml.Unmarshal: %w", err)
// 	}
//
// 	rep := mpd.Periods[0].
// 		AdaptationSets[0].
// 		Representations[0]
//
// 	tmpl := rep.SegmentTemplate
//
// 	segments := []string{tmpl.Initialization}
//
// 	n := tmpl.StartNumber
// 	for _, s := range tmpl.Timeline.Segments {
// 		repeat := max(s.R, 0)
// 		for i := 0; i <= repeat; i++ {
// 			url := strings.Replace(
// 				tmpl.Media,
// 				"$Number$",
// 				strconv.Itoa(n),
// 				1,
// 			)
// 			segments = append(segments, url)
// 			n++
// 		}
// 	}
//
// 	var readers []io.ReadCloser
//
// 	for _, url := range segments {
// 		resp, err := utils.Fetch(ctx, url)
// 		if err != nil {
// 			return nil, fmt.Errorf("fetchAlbum: %w", err)
// 		}
//
// 		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
// 			defer resp.Body.Close()
// 			return nil, fmt.Errorf("fetchAlbum: http: %w", errors.New(resp.Status))
// 		}
//
// 		readers = append(readers, resp.Body)
// 	}
//
// 	fullReader := utils.MultiReadCloser(readers...)
//
// 	return fullReader, nil
// }
//
// func remuxM4AtoFLAC(reader io.ReadCloser) (io.ReadCloser, error) {
// 	cmd := exec.Command("ffmpeg",
// 		"-nostdin",
// 		"-fflags", "+genpts",
// 		"-i", "pipe:0",
// 		"-map", "0:a:0",
// 		"-map_metadata", "0",
// 		"-c:a", "copy",
// 		"-f", "flac",
// 		"pipe:1")
// 	stdin, err := cmd.StdinPipe()
// 	if err != nil {
// 		_ = reader.Close()
// 		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
// 	}
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		_ = reader.Close()
// 		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
// 	}
//
// 	if err := cmd.Start(); err != nil {
// 		_ = reader.Close()
// 		return nil, fmt.Errorf("remuxM4AtoFLAC: %w", err)
// 	}
//
// 	go func() {
// 		if _, err := io.Copy(stdin, reader); err != nil {
// 			log.Println("io.Copy(stdin, reader): %w", err)
// 		}
// 		_ = stdin.Close()
// 		_ = reader.Close()
// 	}()
//
// 	newReader, writer := io.Pipe()
//
// 	go func() {
// 		if _, err := io.Copy(writer, stdout); err != nil {
// 			writer.CloseWithError(err)
// 		}
//
// 		if err := cmd.Wait(); err != nil {
// 			_ = writer.CloseWithError(err)
// 		} else {
// 			_ = writer.Close()
// 		}
// 	}()
// 	return newReader, nil
// }
//
// func (p *Hifi) Download(ctx context.Context, userId uint, id string) (io.ReadCloser, string, error) {
// 	quality := "HI_RES_LOSSLESS"
// 	if user, err := repository.GetUserByID(userId); err != nil {
// 		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
// 	} else if !user.HiRes {
// 		quality = "LOSSLESS"
// 	}
//
// 	instances, err := repository.ListInstancesByUserIDByAPI(userId, p.Name())
// 	if err != nil {
// 		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
// 	}
//
// 	info, err := getDownloadInfo(ctx, instances, id, quality)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
// 	}
//
// 	manifest, err := base64.StdEncoding.DecodeString(info.Data.Manifest)
// 	if err != nil {
// 		return nil, "", fmt.Errorf("Hifi.Download: base64.StdEncoding.DecodeString: %w", err)
// 	}
//
// 	var reader io.ReadCloser
// 	switch info.Data.ManifestMimeType {
// 	case "application/vnd.tidal.bts":
// 		reader, err = downloadTidal(ctx, manifest)
// 	case "application/dash+xml":
// 		reader, err = downloadMPD(ctx, manifest)
// 	default:
// 		err = errors.New("manifest type unknown")
// 	}
// 	if err != nil {
// 		return nil, "", fmt.Errorf("Hifi.Download: %w", err)
// 	}
//
// 	var extension string
// 	switch info.Data.AudioQuality {
// 	case "HI_RES_LOSSLESS":
// 		if quality != "HI_RES_LOSSLESS" {
// 			return nil, "", fmt.Errorf("Hifi.Download: %w", errors.New("audio quality received not conform"))
// 		}
// 		reader, err = remuxM4AtoFLAC(reader)
// 		if err != nil {
// 			return nil, "", fmt.Errorf("Hifi.Download: %w", err)
// 		}
// 		extension = "flac"
// 	case "LOSSLESS":
// 		extension = "flac"
// 	case "HIGH":
// 		extension = "m4a"
// 	}
//
// 	return reader, extension, nil
// }
