package hifiv2_2

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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Hifi.Status: http.DefaultClient.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Hifi.Status: %w", errors.New("Bad Status code"))
	}

	var status status
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return fmt.Errorf("Hifi.Status: json.Decode: %w", err)
	}

	if !(status.Version == "2.2" && status.Repo == "https://github.com/uimaxbai/hifi-api") {
		return fmt.Errorf("Hifi.Status: %w", errors.New("Status content don't match"))
	}

	return nil
}
