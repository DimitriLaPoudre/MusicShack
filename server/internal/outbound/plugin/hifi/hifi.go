package hifi

// https://github.com/uimaxbai/hifi-api

import (
	"context"
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

func (p *Hifi) Lyrics(ctx context.Context, instances []model.Instance, id string) (string, string, error) {
	return "", "", nil
}
