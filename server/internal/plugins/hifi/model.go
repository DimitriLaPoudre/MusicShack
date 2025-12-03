package hifi

import "context"

type Hifi struct{}

func (p *Hifi) Name() string {
	return "hifi"
}

func (p *Hifi) Cover(ctx context.Context, id string) (string, error) {
	return "", nil
}
func (p *Hifi) Lyrics(ctx context.Context, id string) (string, string, error) {
	return "", "", nil
}
