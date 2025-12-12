package hifiv2

// https://github.com/uimaxbai/hifi-api

import (
	"context"
)

type HifiV2 struct{}

func (p *HifiV2) Name() string {
	return "hifiV2"
}

func (p *HifiV2) Cover(ctx context.Context, id string) (string, error) {
	return "", nil
}
func (p *HifiV2) Lyrics(ctx context.Context, id string) (string, string, error) {
	return "", "", nil
}
