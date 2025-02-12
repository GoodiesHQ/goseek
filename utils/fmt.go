package utils

import "strings"

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
	var builder strings.Builder

	builder.WriteString(`<a href="`)
	builder.WriteString(Format(path, trail))
	if apikey != "" {
		builder.WriteString("?apikey=")
		builder.WriteString(apikey)
	}
	builder.WriteString(`">`)
	builder.WriteString(strings.TrimPrefix(Format(name, trail), "/"))
	builder.WriteString("<br />\n")

	return builder.String()
}
