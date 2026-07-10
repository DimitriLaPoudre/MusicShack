package hifi_utils

import "github.com/Ascension-EIP/Ascension/apps/server/internal/model"

func InstancesToUrls(instances []model.Instance) []string {
	urls := []string{}
	for _, i := range instances {
		urls = append(urls, i.Url)
	}
	return urls
}
