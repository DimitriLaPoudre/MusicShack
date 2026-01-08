package hifiv2_2

// https://github.com/uimaxbai/hifi-api

import "context"

type Hifi struct{}

func (p *Hifi) Name() string {
	return "hifiV2.2"
}

func (p *Hifi) Provider() string {
	return "tidal"
}

func (p *Hifi) Priority() int {
	return 1
}

func (p *Hifi) Lyrics(ctx context.Context, userId uint, id string) (string, string, error) {
	return "", "", nil
}
