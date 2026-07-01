package usecase

import (
	"context"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type InstanceUseCase struct {
	l    *zerolog.Logger
	repo model.InstanceRepository
}

func NewInstanceUseCase(l *zerolog.Logger, repo model.InstanceRepository) InstanceUseCase {
	return InstanceUseCase{
		l:    l,
		repo: repo,
	}
}

func (s *InstanceUseCase) CreateInstance(c context.Context, i *model.Instance) (*model.Instance, error) {
	return s.repo.CreateInstance(c, i)
}

func (s *InstanceUseCase) ListInstancesByUserID(c context.Context, userID uuid.UUID) ([]*model.Instance, error) {
	return s.repo.ListInstancesByUserID(c, userID)
}

func (s *InstanceUseCase) DeleteInstanceByUserID(c context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repo.DeleteInstanceByUserID(c, id, userID)
}
