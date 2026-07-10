package hifi

// https://github.com/uimaxbai/hifi-api

import (
	"context"
	"io"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Hifi struct {
	l        *zerolog.Logger
	limiters map[string]*rate.Limiter
}

func NewHifi(l *zerolog.Logger) Hifi {
	return Hifi{
		l:        l,
		limiters: map[string]*rate.Limiter{},
	}
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
