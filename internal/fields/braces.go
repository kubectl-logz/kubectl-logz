package fields

import (
	"regexp"
	"strings"
)

var braceRx = regexp.MustCompile("\\[([^]]*)]")

var Braces = func(l string) []string {
	find := braceRx.FindAllStringSubmatch(l, -1)
	chunk := make([]string, len(find))
	for i, found := range find {
		chunk[i] = found[1]
	}
	return append(chunk, strings.Trim(braceRx.ReplaceAllString(l, ""), " -"))
}
