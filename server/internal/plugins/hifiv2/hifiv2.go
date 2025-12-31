package hifiv2

// https://github.com/uimaxbai/hifi-api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HifiV2 struct{}

func (p *HifiV2) Name() string {
	return "hifiV2"
}

func (p *HifiV2) Status(ctx context.Context, url string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("HifiV2.Status: http.NewRequestWithContext: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("HifiV2.Status: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HifiV2.Status: %w", errors.New("Bad Status code"))
	}

	var status status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return fmt.Errorf("HifiV2.Status: json.Decode: %w", err)
	}

	if status.Version != "2.0" || status.HifiApi != "v2.0" || status.Repo != "https://github.com/uimaxbai/hifi-api" {
		return fmt.Errorf("HifiV2.Status: %w", errors.New("Status content don't match"))
	}

	return nil
}

func (p *HifiV2) Lyrics(ctx context.Context, userId uint, id string) (string, string, error) {
	return "", "", nil
}
