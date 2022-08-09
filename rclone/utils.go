package rclone

import (
	"strings"
)

// Remove '\n' from a string.
func removeNewlines(s string) string {
	return strings.Replace(s, "\n", " ", -1)
}

// Add a param e.g. "patato=1" or "tomato", to url.
func addParamToURL(url, param string) string {
	questionMark := false
	for _, c := range url {
		if c == '/' {
			questionMark = false
		} else if c == '?' {
			questionMark = true
		}
	}
	if questionMark {
		url = url + "&" + param
	} else {
		url = url + "?" + param
	}
	return url
}
