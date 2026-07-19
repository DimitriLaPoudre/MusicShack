package hifi

// https://github.com/uimaxbai/hifi-api

import (
	"context"
	"io"
	"sync"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type Hifi struct {
	limiters sync.Map // map[string]*rate.Limiter
}

func NewHifi() Hifi {
	return Hifi{}
}

func (p *Hifi) Name() string {
	return "hifi"
}

func (p *Hifi) Provider() string {
	return "tidal"
}

func (p *Hifi) Download(ctx context.Context, user *model.User, id string) (io.ReadCloser, string, error) {
	return nil, "", nil
}

func (p *Hifi) Url(ctx context.Context, instances []model.Instance, url string) (model.UrlItem, error) {
	return model.UrlItem{}, nil
}

func (p *Hifi) Lyrics(ctx context.Context, instances []model.Instance, id string) (string, string, error) {
	return "", "", nil
}
