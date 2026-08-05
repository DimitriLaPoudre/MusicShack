package hifi

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	hifi_utils "github.com/DimitriLaPoudre/MusicShack/internal/outbound/plugin/hifi/utils"
)

func (p *Hifi) getSongInfo(ctx context.Context, urls []string, id string) (songResponse, error) {
	songInfo, err := hifi_utils.MultiFetchTyped[songResponse](ctx, urls, "/info/?id="+url.QueryEscape(id), &p.limiters)
	if err != nil {
		return songResponse{}, fmt.Errorf("fetch song info with url list: %w", err)
	}

	return songInfo, nil
}

func (p *Hifi) getSongDownloadInfo(ctx context.Context, urls []string, id string) (downloadResponse, error) {
	downloadInfo, err := hifi_utils.MultiFetchTyped[downloadResponse](ctx, urls, "/track/?id="+url.QueryEscape(id), &p.limiters)
	if err != nil {
		return downloadResponse{}, fmt.Errorf("fetch song download info with url list: %w", err)
	}

	return downloadInfo, nil
}

func (p *Hifi) getSong(ctx context.Context, urls []string, id string) (songResponse, downloadResponse, error) {
	var songInfo songResponse
	var songInfoErr error
	var downloadInfo downloadResponse
	// var downloadInfoErr error
	var wg sync.WaitGroup

	wg.Go(func() {
		songInfo, songInfoErr = p.getSongInfo(ctx, urls, id)
	})
	wg.Go(func() {
		downloadInfo, _ = p.getSongDownloadInfo(ctx, urls, id)
	})
	wg.Wait()

	if songInfoErr != nil {
		return songResponse{}, downloadResponse{}, songInfoErr
	}

	return songInfo, downloadInfo, nil
}

func (p *Hifi) SongInfo(ctx context.Context, instances []model.Instance, id string) (model.SongInfo, error) {
	urls := hifi_utils.InstancesToURLs(instances)

	songInfo, downloadInfo, err := p.getSong(ctx, urls, id)
	if err != nil {
		return model.SongInfo{}, err
	}

	song := songInfo.Data.ToSongInfo(1280)
	song.ReplayGain = downloadInfo.Data.TrackReplayGain
	song.AlbumPeak = downloadInfo.Data.AlbumPeakAmplitude

	return song, nil
}
