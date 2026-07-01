package usecase

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/service"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type InstanceUseCase struct {
	l        *zerolog.Logger
	cfg      config.LibraryConfig
	instance service.InstanceService
}

func NewInstanceUseCase(l *zerolog.Logger, cfg config.LibraryConfig, instance service.InstanceService) InstanceUseCase {
	return InstanceUseCase{
		l:        l,
		cfg:      cfg,
		instance: instance,
	}
}

func (s *InstanceUseCase) CreateInstance(c context.Context, i *model.Instance, userID uuid.UUID) (*model.Instance, error) {
	return s.instance.CreateInstance(c, i, userID)
}

func (s *InstanceUseCase) GetInstanceByID(c context.Context, id uuid.UUID) (*model.Instance, error) {
	return s.instance.GetInstanceByID(c, id)
}

func (s *InstanceUseCase) ListInstancesByUserID(c context.Context, userID uuid.UUID) ([]*model.Instance, error) {
	return s.instance.ListInstancesByUserID(c, userID)
}

func (s *InstanceUseCase) DeleteInstance(c context.Context, id uuid.UUID) error {
	return s.instance.DeleteInstance(c, id)
}
