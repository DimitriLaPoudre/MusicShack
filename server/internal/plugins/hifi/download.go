package hifi

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/services"
)

func (p *Hifi) DownloadSong(ctx context.Context, userId uint, id string, quality string) error {
	services.DownloadManager.Add(userId, p, id, quality)
	return nil
}

func (p *Hifi) DownloadAlbum(ctx context.Context, userId uint, id string, quality string) error {
	return nil
}

func (p *Hifi) DownloadArtist(ctx context.Context, userId uint, id string, quality string) error {
	return nil
}
