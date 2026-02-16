package plugins

// https://github.com/binimum/hifi-api

import hifi "github.com/DimitriLaPoudre/MusicShack/server/internal/plugins/hifi"

func init() {
	Register(&hifi.Hifi{})
}
