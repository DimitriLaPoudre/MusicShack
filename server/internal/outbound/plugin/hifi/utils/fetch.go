package hifi_utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/network"
	"golang.org/x/time/rate"
)

func FetchType[T any](ctx context.Context, url string, limiter *rate.Limiter) (T, error) {
	var zero T

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := Fetch(ctx, url, limiter)
	if err != nil {
		return zero, fmt.Errorf("hifi_utils.FetchType: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, fmt.Errorf("hifi_utils.FetchType: http: %w", errors.New(resp.Status))
	}

	var data T
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return zero, fmt.Errorf("hifi_utils.FetchType: json.Decode: %w", err)
	}

	return data, nil
}

func Fetch(ctx context.Context, url string, limiter *rate.Limiter) (*http.Response, error) {
	if !limiter.Allow() {
		return nil, fmt.Errorf("hifi_utils.Fetch: %w", model.ErrPluginRateLimit)
	}
	return network.Fetch(ctx, url)
}
