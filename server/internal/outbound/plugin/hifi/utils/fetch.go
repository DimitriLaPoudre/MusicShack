package hifi_utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/network"
	"github.com/DimitriLaPoudre/MusicShack/internal/pkg/sync"
	"golang.org/x/time/rate"
)

func MultiFetchTyped[T any](ctx context.Context, urls []string, path string, limiters *sync.Map[string, *rate.Limiter]) (T, error) {
	var data T
	err := model.ErrPluginNoInstances
	for _, url := range urls {
		limiter, _ := limiters.LoadOrStore(url, rate.NewLimiter(rate.Every(6*time.Second), 150))

		data, err = FetchTyped[T](ctx, url, path, limiter)
		if err == nil {
			break
		}
	}
	return data, err
}

func FetchTyped[T any](ctx context.Context, url string, path string, limiter *rate.Limiter) (T, error) {
	var zero T

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := Fetch(ctx, url+path, limiter)
	if err != nil {
		return zero, fmt.Errorf("fetch url %s: %w", url+path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, fmt.Errorf("fetch url %s response status: %w", url+path, errors.New(resp.Status))
	}

	var data T
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return zero, fmt.Errorf("decode url %s response: %w", url+path, err)
	}

	return data, nil
}

func Fetch(ctx context.Context, url string, limiter *rate.Limiter) (*http.Response, error) {
	if !limiter.Allow() {
		return nil, model.ErrPluginRateLimit
	}
	return network.Fetch(ctx, url)
}
