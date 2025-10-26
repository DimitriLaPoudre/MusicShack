package plugins

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/models"
)

type Plugin interface {
	Name() string
	Download(string, string) error
	Song(string) (models.SongData, error)
	Album(string) (models.AlbumData, error)
	Artist(string) (models.ArtistData, error)
	Search(string, string, string) (any, error)
	Cover(string) (string, error)
	Lyrics(string) (string, string, error)
}

var registry = make(map[string]Plugin)

func Register(p Plugin) {
	registry[p.Name()] = p
}

func Get(name string) (Plugin, bool) {
	p, ok := registry[name]
	return p, ok
}
