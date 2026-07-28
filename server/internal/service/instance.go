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

func (s *InstanceService) CreateInstance(ctx context.Context, i model.Instance) (model.Instance, error) {
	ping_start := time.Now()
	plugin, err := s.plugin.GetOriginalPlugin(ctx, i.URL)
	if err != nil {
		return model.Instance{}, fmt.Errorf("get plugin design for url %s: %w", i.URL, err)
	}
	ping := time.Since(ping_start)

	i.Plugin = plugin.Name()
	i.Provider = plugin.Provider()

	id, err := uuid.NewV7()
	if err != nil {
		return model.Instance{}, fmt.Errorf("create id for new instance: %w", err)
	}
	i.ID = id
	i.Ping = &ping

	i, err = s.repo.CreateInstance(ctx, i)
	if err != nil {
		return model.Instance{}, fmt.Errorf("create new instance: %w", err)
	}

	return i, nil
}

func (s *InstanceService) RefreshListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Instance, error) {
	instances, err := s.repo.ListInstancesByFilter(ctx, model.InstanceFilter{UserID: &userID})
	if err != nil {
		return []model.Instance{}, fmt.Errorf("list instance for user %s: %w", userID.String(), err)
	}

	for i, instance := range instances {
		ping_start := time.Now()
		err := s.plugin.GetStatus(ctx, instance.Plugin, instance.URL)
		ping := time.Since(ping_start)
		if err != nil {
			_, err = s.repo.UpdateInstancePing(ctx, instance.ID, nil)
			instances[i].Ping = nil
		} else {
			_, err = s.repo.UpdateInstancePing(ctx, instance.ID, &ping)
			instances[i].Ping = &ping
		}
	}

	return instances, nil
}

func (s *InstanceService) ListInstancesByFilter(ctx context.Context, filter model.InstanceFilter) ([]model.Instance, error) {
	instances, err := s.repo.ListInstancesByFilter(ctx, filter)
	if err != nil {
		return []model.Instance{}, fmt.Errorf("list instance for filter %v: %w", filter, err)
	}

	return instances, nil
}

func (s *InstanceService) DeleteInstanceByUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if err := s.repo.DeleteInstanceByUserID(ctx, id, userID); err != nil {
		return fmt.Errorf("delete instance %s of user %s: %w", id.String(), userID.String(), err)
	}

	return nil
}
