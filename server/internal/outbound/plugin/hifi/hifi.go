package hifi

// https://github.com/binimum/hifi-api

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/sync"
	"golang.org/x/time/rate"
)

type Hifi struct {
	limiters sync.Map[string, *rate.Limiter]
}

func NewHifi() Hifi {
	return Hifi{
		limiters: sync.NewMap[string, *rate.Limiter](),
	}
}

func (p *Hifi) Name() string {
	return "hifi"
}

func (p *Hifi) Provider() string {
	return "tidal"
}
