package plugins

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

var registry = make(map[string]models.Plugin)

func Register(p models.Plugin) {
	registry[p.Name()] = p
}

func Get(name string) (models.Plugin, bool) {
	p, ok := registry[name]
	return p, ok
}

func GetRegistry() map[string]models.Plugin {
	return registry
}
