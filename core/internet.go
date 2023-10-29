package core

import (
	"net/url"
)

// takes a fully qualified url and returns the domain name
func GetDomainFromUrl(fqdnURL string) (string, error) {
	parsed, err := url.ParseRequestURI(fqdnURL)
	if AssertErrorToNilf("could not parse url: %w", err) {
		return "", err
	}

	return parsed.Hostname(), nil
}
