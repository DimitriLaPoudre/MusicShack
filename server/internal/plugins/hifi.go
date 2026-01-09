package plugins

// https://github.com/uimaxbai/hifi-api

import hifi "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins/hifi"

func init() {
	Register(&hifi.Hifi{})
}
