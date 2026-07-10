package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/pkg/network"
)

func (p *Hifi) Status(ctx context.Context, url string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := network.Fetch(ctx, url)
	if err != nil {
		return fmt.Errorf("Hifi.Status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Hifi.Status: http: %w", errors.New(resp.Status))
	}

	var status statusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return fmt.Errorf("Hifi.Status: json.Decode: %w", err)
	}

	return nil
}
