package plugins

import (
	"context"

	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

type Plugin interface {
	Name() string
	Download(context.Context, string, string) error
	Song(context.Context, string) (models.SongData, error)
	Album(context.Context, string) (models.AlbumData, error)
	Artist(context.Context, string) (models.ArtistData, error)
	Search(context.Context, string, string, string) (models.SearchData, error)
	Cover(context.Context, string) (string, error)
	Lyrics(context.Context, string) (string, string, error)
}

var registry = make(map[string]Plugin)

func Register(p Plugin) {
	registry[p.Name()] = p
}

func Get(name string) (Plugin, bool) {
	p, ok := registry[name]
	return p, ok
}

func GetRegistry() map[string]Plugin {
	return registry
}
