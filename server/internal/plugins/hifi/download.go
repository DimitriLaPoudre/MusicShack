package hifi

import (
	"context"
	"io"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func (p *Hifi) Download(ctx context.Context, userId uint, id string, quality string, data chan<- models.SongData) (io.ReadCloser, string, error) {
	return nil, "", nil
}
