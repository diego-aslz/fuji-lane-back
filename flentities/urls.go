package flentities

import (
	"net/url"
	"regexp"
	"strings"
)

func generateSlug(name string) string {
	// Replace one or more spaces by a dash
	name = regexp.MustCompile("\\s+").ReplaceAllString(strings.TrimSpace(name), "-")
	return url.QueryEscape(strings.ToLower(name))
}
