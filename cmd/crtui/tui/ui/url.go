package ui

import "strings"

func TrimURLScheme(url string) string {
	uri := strings.TrimPrefix(url, "https://")
	uri = strings.TrimPrefix(uri, "http://")
	return uri
}
