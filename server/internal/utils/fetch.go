package utils

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/sync/semaphore"
)

var limit *semaphore.Weighted = semaphore.NewWeighted(50)

func Fetch(ctx context.Context, url string) (*http.Response, error) {
	if err := limit.Acquire(ctx, 1); err != nil {
		return nil, fmt.Errorf("utils.Fetch: limit.Acquire: %w", err)
	}
	defer limit.Release(1)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("utils.Fetch: http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("utils.Fetch: http.DefaultClient.Do: %w", err)
	}
	return resp, nil
}
