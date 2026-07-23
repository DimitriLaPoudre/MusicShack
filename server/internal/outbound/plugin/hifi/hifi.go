package hifi

// https://github.com/binimum/hifi-api

import (
	"context"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/sync"
	"golang.org/x/time/rate"
)

type Hifi struct {
	limiters sync.Map[string, *rate.Limiter]
	cache    sync.Cache[string]
}

func NewHifi() Hifi {
	return Hifi{
		limiters: sync.NewMap[string, *rate.Limiter](),
		cache:    sync.NewCache[string](20 * time.Minute),
	}
}

func (p *Hifi) Name() string {
	return "hifi"
}

func (p *Hifi) Provider() string {
	return "tidal"
}

func (p *Hifi) Lyrics(ctx context.Context, instances []model.Instance, id string) (string, string, error) {
	return "", "", nil
}
