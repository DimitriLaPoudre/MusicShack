package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (p *Hifi) Status(ctx context.Context, url string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("Hifi.Status: http.NewRequestWithContext: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Hifi.Status: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Hifi.Status: http: %w", errors.New(resp.Status))
	}

	var status status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return fmt.Errorf("Hifi.Status: json.Decode: %w", err)
	}

	if status.Version != "2.2" || status.Repo != "https://github.com/uimaxbai/hifi-api" {
		return fmt.Errorf("Hifi.Status: %w", errors.New("status content don't match"))
	}

	return nil
}
