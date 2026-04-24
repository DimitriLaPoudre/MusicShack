package hifi_utils

import (
	"fmt"
	"strings"
)

func GetImageURL(encryptPath string, size uint) string {
	if encryptPath == "" {
		return ""
	}
	return fmt.Sprintf("https://resources.tidal.com/images/%s/%dx%d.jpg", strings.ReplaceAll(encryptPath, "-", "/"), size, size)
}
