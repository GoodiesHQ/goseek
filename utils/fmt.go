package utils

import "strings"

const APIKEY_PREFIX_LEN = 7

func Format(path string, trail bool) string {
	if path == "" || path == "." {
		return "/"
	}
	result := "/" + strings.Trim(path, "/\\")
	if trail {
		result = result + "/"
	}
	return result
}

func Href(path, name string, trail bool, apikey string) string {
	// string builder for creating an anchor tag
	var builder strings.Builder

	// build the anchor tag
	builder.WriteString(`<a href="`)
	builder.WriteString(Format(path, trail))
	if apikey != "" {
		// add the API key to the URL
		builder.WriteString("?apikey=")
		builder.WriteString(apikey)
	}
	builder.WriteString(`">`)
	builder.WriteString(strings.TrimPrefix(Format(name, trail), "/"))
	builder.WriteString("</a><br />\n")

	return builder.String()
}

func ApiKeyPrefix(apikey string) string {
	// show only the first few bytes of the API key
	return apikey[:APIKEY_PREFIX_LEN] + "..."
}
