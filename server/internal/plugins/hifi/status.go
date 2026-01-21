package hifi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/utils"
)

func (p *Hifi) Status(ctx context.Context, url string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := utils.Fetch(ctx, url)
	if err != nil {
		return fmt.Errorf("Hifi.Status: %w", err)
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
