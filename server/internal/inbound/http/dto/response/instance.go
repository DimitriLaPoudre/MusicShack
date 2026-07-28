package response

import (
	"github.com/DimitriLaPoudre/MusicShack/internal/model"
)

type Instance struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	Provider string `json:"provider"`
	Plugin   string `json:"plugin"`
	URL      string `json:"url"`
	Ping     *int   `json:"ping"`
}

func InstanceToResponse(instance model.Instance) Instance {
	var ping *int
	if instance.Ping == nil {
		ping = nil
	} else {
		tmp := int(instance.Ping.Milliseconds())
		ping = &tmp
	}
	return Instance{
		ID:       instance.ID.String(),
		UserID:   instance.UserID.String(),
		Provider: instance.Provider,
		Plugin:   instance.Plugin,
		URL:      instance.URL,
		Ping:     ping,
	}
}

func InstancesToResponse(instances []model.Instance) []Instance {
	r := []Instance{}
	for _, instance := range instances {
		r = append(r, InstanceToResponse(instance))
	}
	return r
}
