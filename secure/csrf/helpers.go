package csrf

import (
	"net/url"
)

// stringInSlice checks if a string is inside of a slice
func stringInSlice(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

// sameOrigin makes sure that two url.URL values are from the same origin
func sameOrigin(referer, url *url.URL) bool {
	return (referer.Scheme == url.Scheme && referer.Host == url.Host)
}
