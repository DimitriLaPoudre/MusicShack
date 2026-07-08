package hifi

// https://github.com/uimaxbai/hifi-api

import (
	"context"
	"io"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

type Hifi struct {
	l       *zerolog.Logger
	limiter *rate.Limiter
}

func NewHifi(l *zerolog.Logger) Hifi {
	return Hifi{
		l:       l,
		limiter: rate.NewLimiter(rate.Every(6*time.Second), 150),
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

func (p *Hifi) Search(ctx context.Context, url string, song string, album string, artist string) (model.Search, error) {
	return model.Search{}, nil
}

func (p *Hifi) Url(ctx context.Context, url string, id string) (model.UrlItem, error) {
	return model.UrlItem{}, nil
}

func (p *Hifi) Lyrics(ctx context.Context, url string, id string) (string, string, error) {
	return "", "", nil
}
