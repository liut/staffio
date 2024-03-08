package backends

import (
	"net/url"
	"strings"
)

func SiteFromDomain(name string) string {
	if pos := strings.Index(name, "."); pos > 0 {
		return name[:pos]
	}
	return name
}

func SiteFromURL(uri string) string {
	parsed, err := url.ParseRequestURI(uri)
	if err != nil {
		return SiteFromDomain(uri)
	}
	return SiteFromDomain(parsed.Host)
}
