package hifi

import (
	"context"
	"errors"
	"fmt"
)

type Hifi struct{}

func (p *Hifi) Name() string {
	return "hifi"
}

func (p *Hifi) Status(ctx context.Context, url string) error {
	return fmt.Errorf("Hifi.Status: %w", errors.New("Hifi is deprecated"))
}

func (p *Hifi) Lyrics(ctx context.Context, userId uint, id string) (string, string, error) {
	return "", "", nil
}
