package plugins

// https://github.com/sachinsenal0x64/hifi

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugin"
)

type Hifi struct{}

func (p *Hifi) Test() string {
	return "test"
}

func init() {
	plugin.Register("hifi", &Hifi{})
}
