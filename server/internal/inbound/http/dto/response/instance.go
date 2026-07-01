package response

import (
	"time"

	"github.com/Ascension-EIP/Ascension/apps/server/internal/model"
)

type Instance struct {
	ID       string        `json:"id"`
	UserID   string        `json:"user_id"`
	Provider string        `json:"provider"`
	Plugin   string        `json:"plugin"`
	Url      string        `json:"url"`
	Ping     time.Duration `json:"ping"`
}

func InstanceToResponse(instance *model.Instance) *Instance {
	return &Instance{
		ID:       instance.ID.String(),
		UserID:   instance.UserID.String(),
		Provider: instance.Provider,
		Plugin:   instance.Plugin,
		Url:      instance.Url,
		Ping:     instance.Ping,
	}
}

func InstancesToResponse(instances []*model.Instance) []*Instance {
	r := []*Instance{}
	for _, instance := range instances {
		r = append(r, InstanceToResponse(instance))
	}
	return r
}
