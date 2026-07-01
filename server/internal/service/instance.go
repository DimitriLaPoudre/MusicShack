package service

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/Ascension-EIP/Ascension/apps/server/internal/setup/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type InstanceService struct {
	l    *zerolog.Logger
	cfg  config.LibraryConfig
	repo model.InstanceRepository
}

func NewInstanceService(l *zerolog.Logger, cfg config.LibraryConfig, repo model.InstanceRepository) InstanceService {
	return InstanceService{
		l:    l,
		cfg:  cfg,
		repo: repo,
	}
}

func (s *InstanceService) CreateInstance(c context.Context, i *model.Instance, userID uuid.UUID) (*model.Instance, error) {
	if i == nil {
		return nil, model.ErrInvalidInput
	}

	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	i.ID = id

	i.UserID = userID

	return s.repo.CreateInstance(c, i)
}

func (s *InstanceService) GetInstanceByID(c context.Context, id uuid.UUID) (*model.Instance, error) {
	return s.repo.GetInstanceByID(c, id)
}

func (s *InstanceService) ListInstancesByUserID(c context.Context, userID uuid.UUID) ([]*model.Instance, error) {
	return s.repo.ListInstancesByUserID(c, userID)
}

func (s *InstanceService) DeleteInstance(c context.Context, id uuid.UUID) error {
	return s.repo.DeleteInstance(c, id)
}
