package hifiv2

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
)

func (p *HifiV2) DownloadSong(ctx context.Context, userId uint, id string, quality string) error {
	services.DownloadManager.Add(userId, p, id, quality)
	return nil
}

func (p *HifiV2) DownloadAlbum(ctx context.Context, userId uint, id string, quality string) error {
	return nil
}

func (p *HifiV2) DownloadArtist(ctx context.Context, userId uint, id string, quality string) error {
	return nil
}
