package plugins

// https://github.com/sachinsenal0x64/hifi

import (
	"github.com/DimitriLaPoudre/MusicShack/server/internal/plugins/hifiv1"
)

func init() {
	Register(&hifiv1.HifiV1{})
}
