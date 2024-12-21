package pkg

import (
	"net/url"
	"strings"
)

func ShortInfo(ShortUrl string) string {
	extractLastSegment := func(link string) string {
		parsedUrl, err := url.Parse(link)
		if err != nil {
			return link
		}
		pathSegments := strings.Split(parsedUrl.Path, "/")
		if len(pathSegments) > 0 {
			return pathSegments[len(pathSegments)-1]
		}
		return ""
	}
	NameShort := extractLastSegment(ShortUrl)
	return NameShort
}

func OriginalInfo(OriginalURL string) string {
	extractDomain := func(link string) string {
		parsedUrl, err := url.Parse(link)
		if err != nil {
			return link
		}
		return parsedUrl.Hostname()
	}
	NameOrig := extractDomain(OriginalURL)
	return NameOrig
}
