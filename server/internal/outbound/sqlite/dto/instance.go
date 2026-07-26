package dto

import (
	"time"

	"github.com/DimitriLaPoudre/MusicShack/internal/model"
	"github.com/google/uuid"
)

type Instance struct {
	ID       uuid.UUID      `db:"id"`
	UserID   uuid.UUID      `db:"user_id"`
	Provider string         `db:"provider"`
	Plugin   string         `db:"plugin"`
	Url      string         `db:"url"`
	Ping     *time.Duration `db:"ping"`
}

func (i Instance) ToInstance() model.Instance {
	return model.Instance{
		ID:       i.ID,
		UserID:   i.UserID,
		Provider: i.Provider,
		Plugin:   i.Plugin,
		Url:      i.Url,
		Ping:     i.Ping,
	}
}

func InstancesToInstances(dto []Instance) []model.Instance {
	instances := []model.Instance{}
	for _, i := range dto {
		instances = append(instances, model.Instance{
			ID:       i.ID,
			UserID:   i.UserID,
			Provider: i.Provider,
			Plugin:   i.Plugin,
			Url:      i.Url,
			Ping:     i.Ping,
		})
	}
	return instances
}
