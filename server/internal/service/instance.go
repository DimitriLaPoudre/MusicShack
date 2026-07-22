package service

import (
	"context"
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type InstanceService struct {
	plugin *PluginService
	repo   model.InstanceRepository
}

func NewInstanceService(plugin *PluginService, repo model.InstanceRepository) InstanceService {
	return InstanceService{
		plugin: plugin,
		repo:   repo,
	}
}

func (s *InstanceService) CreateInstance(c context.Context, i model.Instance) (model.Instance, error) {
	ping_start := time.Now()
	plugin, err := s.plugin.GetOriginalPlugin(c, i.Url)
	if err != nil {
		return model.Instance{}, fmt.Errorf("get plugin design for url %s: %w", i.Url, err)
	}
	ping := time.Since(ping_start)

	i.Plugin = plugin.Name()
	i.Provider = plugin.Provider()

	id, err := uuid.NewV7()
	if err != nil {
		return model.Instance{}, fmt.Errorf("create id for new instance: %w", err)
	}
	i.ID = id
	i.Ping = ping

	i, err = s.repo.CreateInstance(c, i)
	if err != nil {
		return model.Instance{}, fmt.Errorf("create new instance: %w", err)
	}

	return i, nil
}

func (s *InstanceService) ListInstancesByFilter(c context.Context, filter model.InstanceFilter) ([]model.Instance, error) {
	instances, err := s.repo.ListInstancesByFilter(c, filter)
	if err != nil {
		return []model.Instance{}, fmt.Errorf("list instance for filter %v: %w", filter, err)
	}

	return instances, nil
}

func (s *InstanceService) DeleteInstanceByUserID(c context.Context, id uuid.UUID, userID uuid.UUID) error {
	if err := s.repo.DeleteInstanceByUserID(c, id, userID); err != nil {
		return fmt.Errorf("delete instance %s of user %s: %w", id.String(), userID.String(), err)
	}

	return nil
}
