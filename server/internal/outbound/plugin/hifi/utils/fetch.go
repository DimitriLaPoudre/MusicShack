package hifi_utils

import (
	"context"
	"net/http"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/network"
	"golang.org/x/time/rate"
)

func Fetch(ctx context.Context, url string, limiter *rate.Limiter) (*http.Response, error) {
	if !limiter.Allow() {
		return nil, model.ErrPluginRateLimit
	}
	return network.Fetch(ctx, url)
}
