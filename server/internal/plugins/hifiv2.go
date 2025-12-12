package plugins

// https://github.com/uimaxbai/hifi-api

import "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins/hifiv2"

func init() {
	Register(&hifiv2.HifiV2{})
}
