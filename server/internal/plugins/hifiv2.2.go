package plugins

// https://github.com/uimaxbai/hifi-api

import hifiv2_2 "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins/hifiv2.2"

func init() {
	Register(&hifiv2_2.Hifi{})
}
