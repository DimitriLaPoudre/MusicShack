package hifi

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

func (p *Hifi) Download(ctx context.Context, userId uint, id string, quality string, status chan<- models.Status, data chan<- models.SongData) error {
	return nil
}
