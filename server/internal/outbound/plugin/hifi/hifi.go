package hifi

// https://github.com/uimaxbai/hifi-api

import (
	"context"
	"io"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"golang.org/x/time/rate"
)

type Hifi struct {
	limiter *rate.Limiter
}

func NewHifi() Hifi {
	return Hifi{
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

func (p *Hifi) Playlist(ctx context.Context, url string, id string) (model.Playlist, error) {
	return model.Playlist{}, nil
}

func (p *Hifi) Album(ctx context.Context, url string, id string) (model.Album, error) {
	return model.Album{}, nil
}

func (p *Hifi) Artist(ctx context.Context, url string, id string) (model.Artist, error) {
	return model.Artist{}, nil
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
