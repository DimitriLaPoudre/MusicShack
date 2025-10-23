package plugins

// https://github.com/sachinsenal0x64/hifi

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugin"
)

type Hifi struct{}

func (p *Hifi) Name() string {
	return "hifi"
}
func (p *Hifi) Download(id string, quality string) error {
	return nil
}
func (p *Hifi) Song(id string) (any, error) {
	return nil, nil
}
func (p *Hifi) Album(id string) (any, error) {
	return nil, nil
}
func (p *Hifi) Artist(id string) (any, error) {
	return nil, nil
}
func (p *Hifi) Search(song, album, artist string) (any, error) {
	return nil, nil
}
func (p *Hifi) Cover(id string) (string, error) {
	return "", nil
}
func (p *Hifi) Lyrics(id string) (string, string, error) {
	return "", "", nil
}
func init() {
	plugin.Register(&Hifi{})
}
