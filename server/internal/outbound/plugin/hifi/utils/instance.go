package hifi_utils

import "github.com/DimitriLaPoudre/MusicShack/internal/model"

func InstancesToURLs(instances []model.Instance) []string {
	urls := []string{}
	for _, i := range instances {
		urls = append(urls, i.URL)
	}
	return urls
}
