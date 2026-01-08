package hifiv2

// https://github.com/uimaxbai/hifi-api

import "context"

type HifiV2 struct{}

func (p *HifiV2) Name() string {
	return "hifiV2"
}

func (p *HifiV2) Provider() string {
	return "tidal"
}

func (p *HifiV2) Priority() int {
	return 0
}

func (p *HifiV2) Lyrics(ctx context.Context, userId uint, id string) (string, string, error) {
	return "", "", nil
}
