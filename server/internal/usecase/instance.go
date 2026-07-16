package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/DimitriLaPoudre/MusicShack/internal/service"
	"github.com/google/uuid"
)

type InstanceUseCase struct {
	plugin *service.PluginService
	repo   model.InstanceRepository
}

func NewInstanceUseCase(plugin *service.PluginService, repo model.InstanceRepository) InstanceUseCase {
	return InstanceUseCase{
		plugin: plugin,
		repo:   repo,
	}
}

func (s *InstanceUseCase) CreateInstance(c context.Context, i model.Instance) (model.Instance, error) {
	ping_start := time.Now()
	plugin, err := s.plugin.GetOriginalPlugin(c, i.Url)
	if err != nil {
		return model.Instance{}, err
	}
	ping := time.Since(ping_start)

	i.Plugin = plugin.Name()
	i.Provider = plugin.Provider()

	id, err := uuid.NewV7()
	if err != nil {
		return model.Instance{}, fmt.Errorf("InstanceUseCase.CreateInstance: uuid.NewV7 %w", err)
	}
	i.ID = id
	i.Ping = ping

	return s.repo.CreateInstance(c, i)
}

func (s *InstanceUseCase) ListInstancesByFilter(c context.Context, filter model.InstanceFilter) ([]model.Instance, error) {
	return s.repo.ListInstancesByFilter(c, filter)
}

func (s *InstanceUseCase) DeleteInstanceByUserID(c context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.DeleteInstanceByUserID(c, id, userID)
}
