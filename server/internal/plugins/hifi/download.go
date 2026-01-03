package hifi

import (
	"context"
	"io"
)

func (p *Hifi) Download(ctx context.Context, userId uint, id string, quality string) (io.ReadCloser, string, error) {
	return nil, "", nil
}
